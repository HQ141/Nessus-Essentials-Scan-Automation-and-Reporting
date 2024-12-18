package services

import (
	"NessusEssAutomation/config"
	"encoding/base64"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// sendgridMail sends the weekly vulnerability assessment report via email
// using the SendGrid API. It accepts the path to the PDF report as input.
func sendgridMail(report string) error {
	// Load SMTP configuration from the config file or environment variables
	config, err := config.LoadSMTPConfig()
	if err != nil {
		log.Println("Failed to load SMTP configuration")
		return err
	}

	// Create a new SendGrid mail object
	m := mail.NewV3Mail()

	// Set email content as HTML
	content := mail.NewContent("text/html", "<p>Weekly Vulnerability Assessment Report.</p>")
	m.SetFrom(mail.NewEmail("no-reply@vidizmo.com", config.Username))
	m.AddContent(content)

	// Create a personalization object for recipients
	personalization := mail.NewPersonalization()

	// Load email recipients from the specified recipients file
	recievers, err := getReceivers(config.Recipients)
	if err != nil {
		log.Println("Failed to load email recipients")
		return err
	}

	// Add recipients to the personalization object
	for _, _emailaddress := range recievers {
		to := mail.NewEmail("Recipient", _emailaddress)
		personalization.AddTos(to)
	}

	// Set email subject
	personalization.Subject = "Vulnerability Assessment Report"
	m.AddPersonalizations(personalization)

	// Attach the PDF report
	pdf := mail.NewAttachment()

	// Read the PDF file
	dat, err := os.ReadFile(report)
	if err != nil {
		log.Println("Failed to read the report file:", err)
		return err
	}

	// Encode the PDF file content in base64
	encoded := base64.StdEncoding.EncodeToString(dat)
	pdf.SetContent(encoded)
	pdf.SetType("application/pdf")
	pdf.SetFilename(report)
	pdf.SetDisposition("attachment")
	m.AddAttachment(pdf)

	// Create a SendGrid API request
	request := sendgrid.GetRequest(config.Password, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	// Send the email via the SendGrid API
	_, err = sendgrid.API(request)
	if err != nil {
		log.Println("Failed to send mail via SendGrid:", err)
		return err
	}

	log.Println("Mail sent successfully via SendGrid")
	return nil
}
