package mail

// mocks a mail service and allow to access the receive messages
type MockMailService struct {
	SendMails []MailAttributes
}

func (service *MockMailService) SendNotification(attributes MailAttributes) error {
	if service.SendMails == nil {
		service.SendMails = make([]MailAttributes, 0)
	}
	service.SendMails = append(service.SendMails, attributes)

	return nil
}
