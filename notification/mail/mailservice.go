package mail

// MailAttributes contains E-Mail attributes
type MailAttributes struct {
	To      []string
	Subject string
	Content string
}

// Service forwards E-Mail to a set of receivers
type MailService interface {
	AddMessage(attributes MailAttributes)
	SendNotifications() error
}
