package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	recipientEmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	templateNameRegex   = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
)

// ErrInvalidArgument is returned when SendTemplatedEmail receives invalid input.
var ErrInvalidArgument = errors.New("invalid argument")

// ErrQueueFull is returned when the worker pool channel buffer is full.
var ErrQueueFull = errors.New("email queue is full")

type emailTask struct {
	To           string
	TemplateName string
	TemplateData map[string]interface{}
}

// EmailService implements email sending functionality via AWS SES with an async worker pool.
type EmailService struct {
	sesClient        *ses.SES
	cwMetrics        *cloudwatch.CloudWatch
	cwLogs           *cloudwatchlogs.CloudWatchLogs
	metricNamespace  string
	logGroup         string
	logStream        string
	senderEmail      string
	workerPool       chan emailTask
	workerCount      int
	wg               sync.WaitGroup
	shutdownCtx      context.Context
	cancel           context.CancelFunc
	log              *log.Logger
}

// NewEmailService creates a new EmailService with the provided configuration and logger.
func NewEmailService(cfg EmailServiceConfig, logger *log.Logger) (*EmailService, error) {
	creds := credentials.NewStaticCredentials(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, "")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWSRegion),
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	svc := &EmailService{
		sesClient:       ses.New(sess),
		cwMetrics:       cloudwatch.New(sess),
		metricNamespace: cfg.CloudWatchMetricNamespace,
		logGroup:        cfg.CloudWatchLogGroup,
		logStream:       fmt.Sprintf("email-service-%d", time.Now().UnixNano()),
		senderEmail:     cfg.SESSenderEmail,
		workerPool:      make(chan emailTask, cfg.WorkerQueueCapacity),
		workerCount:     cfg.WorkerPoolSize,
		shutdownCtx:     ctx,
		cancel:          cancel,
		log:             logger,
	}

	if cfg.CloudWatchLogGroup != "" {
		svc.cwLogs = cloudwatchlogs.New(sess)
		if initErr := svc.initCloudWatchLogStream(); initErr != nil {
			logger.Printf(`{"level":"WARN","event_type":"cloudwatch_init_failure","error_message":%q,"timestamp":%q}`,
				initErr.Error(), time.Now().UTC().Format(time.RFC3339))
			svc.cwLogs = nil
		}
	}

	return svc, nil
}

func (es *EmailService) initCloudWatchLogStream() error {
	_, err := es.cwLogs.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(es.logGroup),
	})
	if err != nil {
		if aerr, ok := err.(interface{ Code() string }); !ok || aerr.Code() != "ResourceAlreadyExistsException" {
			return err
		}
	}

	_, err = es.cwLogs.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(es.logGroup),
		LogStreamName: aws.String(es.logStream),
	})
	return err
}

// StartWorkerPool launches the worker goroutines.
func (es *EmailService) StartWorkerPool() {
	for i := 0; i < es.workerCount; i++ {
		es.wg.Add(1)
		go es.sendEmailTask()
	}
}

// StopWorkerPool signals workers to stop and waits for them to finish.
func (es *EmailService) StopWorkerPool() {
	es.cancel()
	close(es.workerPool)
	es.wg.Wait()
}

func (es *EmailService) sendEmailTask() {
	defer es.wg.Done()
	for {
		select {
		case <-es.shutdownCtx.Done():
			return
		case task, ok := <-es.workerPool:
			if !ok {
				return
			}
			es.dispatchToSES(task)
		}
	}
}

func (es *EmailService) emitMetric(name string, errorType string) {
	dims := []*cloudwatch.Dimension{}
	if errorType != "" {
		dims = append(dims, &cloudwatch.Dimension{
			Name:  aws.String("error_type"),
			Value: aws.String(errorType),
		})
	}
	go func() {
		_, err := es.cwMetrics.PutMetricData(&cloudwatch.PutMetricDataInput{
			Namespace: aws.String(es.metricNamespace),
			MetricData: []*cloudwatch.MetricDatum{
				{
					MetricName: aws.String(name),
					Dimensions: dims,
					Value:      aws.Float64(1),
					Unit:       aws.String(cloudwatch.StandardUnitCount),
					Timestamp:  aws.Time(time.Now().UTC()),
				},
			},
		})
		if err != nil {
			es.log.Printf(`{"level":"WARN","event_type":"metric_emit_failure","metric_name":%q,"error_message":%q,"timestamp":%q}`,
				name, err.Error(), time.Now().UTC().Format(time.RFC3339))
		}
	}()
}

func (es *EmailService) writeToCloudWatchLogs(message string) {
	if es.cwLogs == nil {
		return
	}
	go func() {
		_, err := es.cwLogs.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(es.logGroup),
			LogStreamName: aws.String(es.logStream),
			LogEvents: []*cloudwatchlogs.InputLogEvent{
				{
					Message:   aws.String(message),
					Timestamp: aws.Int64(time.Now().UnixMilli()),
				},
			},
		})
		if err != nil {
			es.log.Printf(`{"level":"WARN","event_type":"cloudwatch_log_failure","error_message":%q,"timestamp":%q}`,
				err.Error(), time.Now().UTC().Format(time.RFC3339))
		}
	}()
}

func (es *EmailService) logAndShip(msg string) {
	es.log.Print(msg)
	es.writeToCloudWatchLogs(msg)
}

func (es *EmailService) dispatchToSES(task emailTask) {
	dataBytes, err := json.Marshal(task.TemplateData)
	if err != nil {
		msg := fmt.Sprintf(`{"level":"ERROR","event_type":"email_send_failure","email_to":%q,"template_name":%q,"error_type":"marshal_error","error_message":%q,"timestamp":%q}`,
			task.To, task.TemplateName, err.Error(), time.Now().UTC().Format(time.RFC3339))
		es.logAndShip(msg)
		es.emitMetric("email.send.failure.count", "marshal_error")
		return
	}

	input := &ses.SendTemplatedEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(task.To)},
		},
		Source:       aws.String(es.senderEmail),
		Template:     aws.String(task.TemplateName),
		TemplateData: aws.String(string(dataBytes)),
	}

	out, err := es.sesClient.SendTemplatedEmail(input)
	if err != nil {
		requestID := ""
		if out != nil && out.MessageId != nil {
			requestID = *out.MessageId
		}
		errorType := "ses_error"
		if aerr, ok := err.(interface{ Code() string }); ok {
			errorType = aerr.Code()
		}
		msg := fmt.Sprintf(`{"level":"ERROR","event_type":"email_send_failure","email_to":%q,"template_name":%q,"error_type":%q,"error_message":%q,"aws_request_id":%q,"timestamp":%q}`,
			task.To, task.TemplateName, errorType, err.Error(), requestID, time.Now().UTC().Format(time.RFC3339))
		es.logAndShip(msg)
		es.emitMetric("email.send.failure.count", errorType)
		return
	}

	msg := fmt.Sprintf(`{"level":"INFO","event_type":"email_send_success","email_to":%q,"template_name":%q,"timestamp":%q}`,
		task.To, task.TemplateName, time.Now().UTC().Format(time.RFC3339))
	es.logAndShip(msg)
	es.emitMetric("email.send.success.count", "")
}

// SendTemplatedEmail validates input and enqueues the email task for async processing.
func (es *EmailService) SendTemplatedEmail(to string, templateName string, templateData map[string]interface{}) error {
	if to == "" || len(to) > 254 || !recipientEmailRegex.MatchString(to) {
		msg := fmt.Sprintf(`{"level":"WARN","event_type":"email_validation_failure","email_to":%q,"template_name":%q,"error_message":"invalid recipient email","timestamp":%q}`,
			to, templateName, time.Now().UTC().Format(time.RFC3339))
		es.logAndShip(msg)
		return ErrInvalidArgument
	}

	if templateName == "" || !templateNameRegex.MatchString(templateName) {
		msg := fmt.Sprintf(`{"level":"WARN","event_type":"email_validation_failure","email_to":%q,"template_name":%q,"error_message":"invalid template name","timestamp":%q}`,
			to, templateName, time.Now().UTC().Format(time.RFC3339))
		es.logAndShip(msg)
		return ErrInvalidArgument
	}

	task := emailTask{
		To:           to,
		TemplateName: templateName,
		TemplateData: templateData,
	}

	select {
	case es.workerPool <- task:
		return nil
	default:
		msg := fmt.Sprintf(`{"level":"ERROR","event_type":"email_queue_full","email_to":%q,"template_name":%q,"error_message":"worker pool queue is full","timestamp":%q}`,
			to, templateName, time.Now().UTC().Format(time.RFC3339))
		es.logAndShip(msg)
		es.emitMetric("email.send.failure.count", "queue_full")
		return ErrQueueFull
	}
}

// SendPasswordRecoveryEmail sends a password recovery email using the templated email system.
func (es *EmailService) SendPasswordRecoveryEmail(email string, token string) error {
	return es.SendTemplatedEmail(email, "PasswordRecovery", map[string]interface{}{
		"token": token,
	})
}
