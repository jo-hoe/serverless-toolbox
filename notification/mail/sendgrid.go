package mail

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridConfig contain all attributes to initialize the SendGrid mail service
type SendGridConfig struct {
	APIKey        string
	OriginAddress string
	OriginName    string
}

// SendGridService implements MailService
type SendGridService struct {
	config *SendGridConfig
}

// NewSendGridService creates a SendGridService using an already initializes service
func NewSendGridService(config *SendGridConfig) *SendGridService {
	return &SendGridService{
		config: config,
	}
}

// SendMail mail to one or multiple receivers
func (service *SendGridService) SendMail(attributes MailAttributes) error {
	// create new *SGMailV3
	mailObject := mail.NewV3Mail()

	from := mail.NewEmail(service.config.OriginName, service.config.OriginAddress)
	content := mail.NewContent("text/html", attributes.Content)

	mailObject.SetFrom(from)
	mailObject.AddContent(content)

	// create new *Personalization
	personalization := mail.NewPersonalization()

	personalization.Subject = attributes.Subject
	// populate `personalization` with data
	emails := []*mail.Email{}

	for _, mailAddress := range attributes.To {
		mail, _ := mail.ParseEmail(mailAddress)
		emails = append(emails, mail)
	}

	personalization.AddTos(emails...)

	// add `personalization` to `m`
	mailObject.AddPersonalizations(personalization)

	return service.sendRequest(mailObject)
}

func (service *SendGridService) sendRequest(mailObject *mail.SGMailV3) error {
	request := sendgrid.GetRequest(
		service.config.APIKey,
		"/v3/mail/send",
		"https://api.sendgrid.com",
	)

	request.Method = "POST"
	request.Body = mail.GetRequestBody(mailObject)
	result, err := sendgrid.API(request)

	if result.StatusCode != 202 {
		return fmt.Errorf("SendGrid could not send mail. [%d]: %s", result.StatusCode, result.Body)
	}

	return err
}
