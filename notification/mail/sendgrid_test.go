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

const APIKey = ""
const From = ""
const To = ""

func TestSendMail(t *testing.T) {
	if len(APIKey) == 0 || len(From) == 0 || len(To) == 0 {
		t.Skip("Skipping test because integration test variables are not provided.")
	}

	config := SendGridConfig{
		APIKey:        APIKey,
		OriginAddress: From,
		OriginName:    "Tester",
	}

	sender := NewSendGridService(&config)
	attributes := MailAttributes{
		To:      []string{To},
		Subject: "Hi",
		Content: "Another test mail app",
	}

	err := sender.SendMail(attributes)
	checkError(err, t)

	if sender == nil {
		t.Errorf("Sendgrid not initialized")
	}
}

func getTestConfig() SendGridConfig {
	return SendGridConfig{
		APIKey:        "testkey",
		OriginAddress: "keyaddress",
		OriginName:    "testname",
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
