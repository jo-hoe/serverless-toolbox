package mail

import "testing"

func TestMockMailService_SendNotification(t *testing.T) {
	test := MockMailService{}
	attributes := MailAttributes{
		To:      []string{"a@mail.com"},
		Subject: "test subject",
		Content: "test content",
	}

	err := test.SendNotification(attributes)

	if err != nil {
		t.Errorf("Found error %v", err)
	}

	if attributes.To[0] != test.SendMails[0].To[0] ||
		attributes.Subject != test.SendMails[0].Subject ||
		attributes.Content != test.SendMails[0].Content {
		t.Errorf("Excepted not eq actual. Expected: %+v, Actual %+v", attributes, test.SendMails[0])
	}
}
