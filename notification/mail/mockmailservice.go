package mail

// mocks a mail service and allow to access the receive messages
type MockMailService struct {
	SendMails []MailAttributes
}

func NewMockMailService() *MockMailService {
	return &MockMailService{
		SendMails: make([]MailAttributes, 0),
	}
}

func (service MockMailService) SendNotification(attributes MailAttributes) error {
	service.SendMails = append(service.SendMails, attributes)

	return nil
}
