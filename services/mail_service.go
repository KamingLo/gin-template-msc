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

// SendEmail adalah fungsi umum untuk kirim email apa pun
func SendEmail(toEmail, subject, body string) error {
	host := os.Getenv("MAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD") // App Password Gmail
	fromEmail := os.Getenv("MAIL_FROM_ADDRESS")
	fromName := os.Getenv("MAIL_FROM_NAME")

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(fromEmail, fromName))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, username, password)

	// Konfigurasi SSL untuk Gmail port 465
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// SendRegistrationOTP khusus untuk kirim kode verifikasi
func SendRegistrationOTP(toEmail, otp string) error {
	subject := "Verifikasi Akun - Kode OTP"

	tmplPath := filepath.Join("templates", "otp.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("Gagal parse template: %w", err)
	}

	var body bytes.Buffer
	data := DataOTP{OTP: otp}
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("Gagal execute template: %w", err)
	}

	go func(targetEmail, mailSubject, mailBody string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("[Panic Recovery] Error saat kirim email ke %s: %v\n", targetEmail, r)
			}
		}()

		err := SendEmail(targetEmail, mailSubject, mailBody)
		if err != nil {
			fmt.Printf("[Error Background] Gagal kirim email ke %s: %v\n", targetEmail, err)
		}
	}(toEmail, subject, body.String())

	return nil
}
