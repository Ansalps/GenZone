package helper

import (
	"fmt"
	"net/smtp"
	"strings"
)

//var Time time.Time

func SendOTPEmail(email string, otp string) error {
	// // Choose auth method and set it up
	fmt.Println(email, otp)
	fmt.Println("hi")
	// auth := smtp.PlainAuth("", "john.doe@gmail.com", "extremely_secret_pass", "smtp.gmail.com")
	from := "genzoneapi@gmail.com"
	password := "qqeg rvju rbim wcdk" // TODO: Replace with your email password or use a secure method to fetch it
	to := []string{email}
	subject := "OTP for Signup"
	body := "Your OTP is: " + otp

	// Compose email
	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// SMTP server configuration
	smtpServer := "smtp.gmail.com"
	auth := smtp.PlainAuth("", from, password, smtpServer)
	fmt.Println("hello")
	// Send email
	err := smtp.SendMail(smtpServer+":587", auth, from, to, []byte(msg))

	if err != nil {
		fmt.Println("error in send Mail", err.Error())
		return err
	}
	return nil
}
