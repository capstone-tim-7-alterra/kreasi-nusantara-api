package email

import (
	"os"

	mail "github.com/xhit/go-simple-mail/v2"
)

type EmailUtil interface {
	SendOTP(email string, otp string) error
}

type emailUtil struct{}

func NewEmailUtil() *emailUtil {
	return &emailUtil{}
}

func (e *emailUtil) SendOTP(email string, otp string) error {
	server := mail.NewSMTPClient()
	server.Host = os.Getenv("SMTP_HOST")
	server.Port = 587
	server.Username = os.Getenv("SMTP_USERNAME")
	server.Password = os.Getenv("SMTP_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	emailObj := mail.NewMSG()
	emailObj.SetFrom(os.Getenv("EMAIL_FROM")).AddTo(email).SetSubject("Kreasi Nusantara OTP Verification")
	emailObj.SetBody(mail.TextPlain, "Your OTP code is: "+otp)

	return emailObj.Send(smtpClient)
}
