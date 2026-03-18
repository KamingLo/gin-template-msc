package services

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

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

	// Template HTML sederhana & clean
	body := fmt.Sprintf(`
		<div style="font-family: sans-serif; max-width: 400px; margin: auto; border: 1px solid #e5e7eb; padding: 24px; border-radius: 16px;">
			<h2 style="color: #18181b; margin-bottom: 8px;">Halo!</h2>
			<p style="color: #71717a; font-size: 14px;">Gunakan kode OTP di bawah ini untuk menyelesaikan pendaftaran akun kamu.</p>
			<div style="background: #f4f4f5; padding: 16px; border-radius: 12px; text-align: center; margin: 24px 0;">
				<span style="font-size: 32px; font-weight: bold; color: #000; letter-spacing: 4px;">%s</span>
			</div>
			<p style="font-size: 12px; color: #a1a1aa;">Kode ini berlaku selama 5 menit. Jangan bagikan kepada siapa pun.</p>
		</div>
	`, otp)

	return SendEmail(toEmail, subject, body)
}
