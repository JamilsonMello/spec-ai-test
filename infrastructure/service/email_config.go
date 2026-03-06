package service

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// EmailServiceConfig holds configuration for EmailService.
type EmailServiceConfig struct {
	AWSAccessKeyID           string
	AWSSecretAccessKey       string
	AWSRegion                string
	SESSenderEmail           string
	WorkerPoolSize           int
	WorkerQueueCapacity      int
	CloudWatchLogGroup       string
	CloudWatchMetricNamespace string
}

// LoadEmailServiceConfig reads EmailServiceConfig from environment variables.
func LoadEmailServiceConfig() (EmailServiceConfig, error) {
	cfg := EmailServiceConfig{
		AWSAccessKeyID:            os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:        os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:                 os.Getenv("AWS_REGION"),
		SESSenderEmail:            os.Getenv("SES_SENDER_EMAIL"),
		WorkerPoolSize:            10,
		WorkerQueueCapacity:       100,
		CloudWatchLogGroup:        os.Getenv("CLOUDWATCH_LOG_GROUP"),
		CloudWatchMetricNamespace: os.Getenv("CLOUDWATCH_METRIC_NAMESPACE"),
	}

	if cfg.CloudWatchMetricNamespace == "" {
		cfg.CloudWatchMetricNamespace = "EmailService"
	}

	if v := os.Getenv("EMAIL_WORKER_POOL_SIZE"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return EmailServiceConfig{}, fmt.Errorf("invalid EMAIL_WORKER_POOL_SIZE: %w", err)
		}
		cfg.WorkerPoolSize = n
	}

	if v := os.Getenv("EMAIL_WORKER_QUEUE_CAPACITY"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return EmailServiceConfig{}, fmt.Errorf("invalid EMAIL_WORKER_QUEUE_CAPACITY: %w", err)
		}
		cfg.WorkerQueueCapacity = n
	}

	if cfg.SESSenderEmail == "" {
		return EmailServiceConfig{}, errors.New("SES_SENDER_EMAIL is required")
	}
	if !emailRegex.MatchString(cfg.SESSenderEmail) {
		return EmailServiceConfig{}, errors.New("SES_SENDER_EMAIL is not a valid email address")
	}

	return cfg, nil
}
