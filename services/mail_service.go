package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/gomail.v2"
)

type DataOTP struct {
	OTP string
}

type DataResetPassword struct {
	ResetURL string
}

// SendEmail handles the low-level SMTP connection and delivery
func SendEmail(toEmail, subject, body string) error {
	host := os.Getenv("MAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	fromEmail := os.Getenv("MAIL_FROM_ADDRESS")
	fromName := os.Getenv("MAIL_FROM_NAME")

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(fromEmail, fromName))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, username, password)

	// Bypass TLS verification for common mail providers if needed
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// SendRegistrationOTP prepares the template and triggers the background sender
func SendRegistrationOTP(toEmail, otp string) error {
	subject := "Account Verification - OTP Code"

	tmplPath := filepath.Join("templates", "otp.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	data := DataOTP{OTP: otp}
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Send email asynchronously to keep the API fast
	go func(targetEmail, mailSubject, mailBody string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[Panic Recovery] Error sending email to %s: %v\n", targetEmail, r)
			}
		}()

		err := SendEmail(targetEmail, mailSubject, mailBody)
		if err != nil {
			fmt.Printf("[Background Error] Failed to send email to %s: %v\n", targetEmail, err)
		} else {
			fmt.Printf("[Success] OTP Email sent successfully to %s\n", targetEmail)
		}
	}(toEmail, subject, body.String())

	return nil
}

func SendForgotPasswordLink(toEmail, resetURL string) error {
	subject := "Reset Kata Sandi Anda"

	// Arahkan ke file template baru khusus reset password
	tmplPath := filepath.Join("templates", "forgot_password.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	data := DataResetPassword{ResetURL: resetURL}
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Kirim secara asinkron agar API tetap responsif (Non-blocking)
	go func(targetEmail, mailSubject, mailBody string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[Panic Recovery] Error sending email to %s: %v\n", targetEmail, r)
			}
		}()

		err := SendEmail(targetEmail, mailSubject, mailBody)
		if err != nil {
			fmt.Printf("[Background Error] Failed to send email to %s: %v\n", targetEmail, err)
		}
	}(toEmail, subject, body.String())

	return nil
}
