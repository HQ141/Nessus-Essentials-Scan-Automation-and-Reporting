package services

import (
	"NessusEssAutomation/config"
	"log"

	"github.com/wneessen/go-mail"
)

// func SendGridEmail() {
// 	from := mail.NewEmail("Example User", "test@example.com")
// 	subject := "Sending with Twilio SendGrid is Fun"
// 	to := mail.NewEmail("Example User", "test@example.com")
// 	plainTextContent := "and easy to do anywhere, even with Go"
// 	a_pdf := mail.NewAttachment()
// 	dat, err := os.ReadFile("testing.pdf")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	encoded := base64.StdEncoding.EncodeToString([]byte(dat))
// 	a_pdf.SetContent(encoded)
// 	a_pdf.SetType("application/pdf")
// 	a_pdf.SetFilename("testing.pdf")
// 	a_pdf.SetDisposition("attachment")
// 	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
// 	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
// 	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
// 	response, err := client.Send(message)
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		fmt.Println(response.StatusCode)
// 		fmt.Println(response.Body)
// 		fmt.Println(response.Headers)
// 	}
// }

func sendEmail(Filepath string) error {
	config, err := config.LoadSMTPConfig()
	if err != nil {
		return err
	}
	username := config.Username
	password := config.Password

	client, err := mail.NewClient(config.Url, mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthLogin), mail.WithUsername(username), mail.WithPassword(password))
	if err != nil {
		log.Printf("failed to create mail client: %s\n", err)
		return err
	}
	message := mail.NewMsg()
	if err := message.From("test@#test.com"); err != nil {
		log.Fatalf("failed to set FROM address: %s", err)
	}
	recievers, err := getReceivers(config.Recipients)
	if err != nil {
		return err
	}
	for _, _emailaddress := range recievers {
		if err := message.AddTo(_emailaddress); err != nil {
			log.Fatalf("failed to set TO address: %s", err)
		}
	}
	message.Subject("Vulnerability Assessment Report")
	message.SetBodyString(mail.TypeTextPlain, "This is a test for vulnerability assessment report")
	message.AttachFile(Filepath)
	if err = client.DialAndSend(message); err != nil {
		log.Printf("failed to send mail: %s\n", err)
		return err
	}
	log.Printf("Email Sent Successfully")
	return nil
}
