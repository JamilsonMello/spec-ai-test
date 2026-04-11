package mocks

type EmailSenderMock struct {
	Calls []EmailSendCall
}

type EmailSendCall struct {
	Email string
	Token string
}

func NewEmailSenderMock() *EmailSenderMock {
	return &EmailSenderMock{}
}

func (m *EmailSenderMock) SendPasswordRecoveryEmail(email string, token string) error {
	m.Calls = append(m.Calls, EmailSendCall{Email: email, Token: token})
	return nil
}
