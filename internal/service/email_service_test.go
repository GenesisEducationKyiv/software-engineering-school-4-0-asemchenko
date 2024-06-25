package service

import (
	"currency-notifier/internal/mocks"
	"github.com/golang/mock/gomock"
	"gopkg.in/gomail.v2"
	"reflect"
	"testing"
)

//go:generate mockgen -source=./email_service.go -destination=../mocks/email_service.go -package=mocks

// EmailMatcher is a custom matcher for gomail.Message
type EmailMatcher struct {
	expected *gomail.Message
}

func (m *EmailMatcher) Matches(x interface{}) bool {
	msg, ok := x.(*gomail.Message)
	if !ok {
		return false
	}

	return reflect.DeepEqual(m.expected.GetHeader("From"), msg.GetHeader("From")) &&
		reflect.DeepEqual(m.expected.GetHeader("To"), msg.GetHeader("To")) &&
		reflect.DeepEqual(m.expected.GetHeader("Subject"), msg.GetHeader("Subject"))
}

func (m *EmailMatcher) String() string {
	return "matches expected email message"
}

func EqEmail(expected *gomail.Message) gomock.Matcher {
	return &EmailMatcher{expected: expected}
}

func TestEmailService_SendCurrencyRateEmail_shouldCreateCorrectMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dialer := mocks.NewMockDialer(ctrl)
	emailService := NewWithDialer(dialer, "sender@domain.example")

	m := gomail.NewMessage()
	m.SetHeader("From", "sender@domain.example")
	m.SetHeader("To", "client@domain.example")
	m.SetHeader("Subject", "Currency Rate Alert")
	m.SetBody("text/plain", "The current USD to UAH rate is 8.00.")
	m.AddAlternative("text/html", "<strong>The current USD to UAH rate is 8.00.</strong>")

	dialer.EXPECT().DialAndSend(EqEmail(m)).Return(nil)

	// perform the test
	err := emailService.SendCurrencyRateEmail("client@domain.example", 8.0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
