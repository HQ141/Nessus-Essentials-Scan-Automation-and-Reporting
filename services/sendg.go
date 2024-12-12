package services

import (
	"NessusEssAutomation/config"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendgridMail(report string) error {
	config, err := config.LoadSMTPConfig()
	if err != nil {
		return err
	}
	m := mail.NewV3Mail()
	content := mail.NewContent("text/html", "<p>Weekly Vulnerability Assessment Report.</p>")
	from := mail.NewEmail("Sender", config.Username)
	m.SetFrom(from)
	m.AddContent(content)
	personalization := mail.NewPersonalization()
	recievers, err := getReceivers(config.Recipients)
	if err != nil {
		return err
	}
	for _, _emailaddress := range recievers {
		to := mail.NewEmail("test1", _emailaddress)
		personalization.AddTos(to)
	}
	personalization.Subject = "Vulnerability Assessment Report"
	m.AddPersonalizations(personalization)
	pdf := mail.NewAttachment()
	dat, err := os.ReadFile(report)
	if err != nil {
		fmt.Println(err)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(dat))
	pdf.SetContent(encoded)
	pdf.SetType("application/pdf")
	pdf.SetFilename(report)
	pdf.SetDisposition("attachment")
	m.AddAttachment(pdf)
	request := sendgrid.GetRequest(config.Password, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	_, err = sendgrid.API(request)
	if err != nil {
		log.Println("Failed to Send mail")
		return err
	}
	return nil
}
