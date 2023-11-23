package Email

// import (
// 	"fmt"
// 	"log"
// 	"net/smtp"
// 	"strings"

// 	Utils "notification-server/utils"
// )

// type Mail struct {
// 	Sender  string
// 	To      []string
// 	Subject string
// 	Body    string
// }

// func BuildMessage(mail Mail) string {
// 	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
// 	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
// 	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
// 	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
// 	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

// 	return msg
// }

// func SendEmail(title string, message string, sender string, to []string) {

// 	user := Utils.GetEnv("EMAIL_USER")
// 	password := Utils.GetEnv("EMAIL_PASSWORD")
// 	EMAIL_SMTP_HOST := Utils.GetEnv("EMAIL_SMTP_HOST")
// 	EMAIL_SMTP_PORT := Utils.GetEnv("EMAIL_SMTP_PORT")

// 	subject := "Notification service"
// 	body := `<p>` + message + `</p>`

// 	request := Mail{
// 		Sender:  sender,
// 		To:      to,
// 		Subject: subject,
// 		Body:    body,
// 	}

// 	addr := EMAIL_SMTP_HOST + ":" + EMAIL_SMTP_PORT
// 	host := EMAIL_SMTP_HOST

// 	msg := BuildMessage(request)
// 	auth := smtp.PlainAuth("", user, password, host)
// 	err := smtp.SendMail(addr, auth, sender, to, []byte(msg))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println("Email sent successfully")

// }
