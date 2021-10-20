package mail

import (
	"testing"
)

func Test_Init(t *testing.T) {
	config := getTestConfig()

	sender := NewSendGridService(&config)

	if sender == nil {
		t.Errorf("Sendgrid not initialized")
	}
}

func Test_AddMessage(t *testing.T) {
	config := getTestConfig()

	sender := NewSendGridService(&config)
	sender.AddMessage(MailAttributes{
		To:      []string{"test@test.com"},
		Subject: "test",
		Content: "test content",
	})

	if len(sender.messages) != 1 {
		t.Errorf("Expected to see 1 message but found %d", len(sender.messages))
	}
}
