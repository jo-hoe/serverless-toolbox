package mail

import (
	"testing"
)

func TestInit(t *testing.T) {
	config := getTestConfig()

	sender := NewSendGridService(&config)

	if sender == nil {
		t.Errorf("Sendgrid not initialized")
	}
}