package services

import (
	"bytes"
	"embed"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"text/template"
)

//go:embed emails/templates/*.html
var emailTemplates embed.FS

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
	Enabled   bool
}

// GetEmailConfig loads email configuration from environment variables
func GetEmailConfig() *EmailConfig {
	enabled := os.Getenv("SMTP_ENABLED")
	isEnabled := enabled == "true" || enabled == "1"

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		port = 587 // Default SMTP port
	}

	return &EmailConfig{
		Host:      os.Getenv("SMTP_HOST"),
		Port:      port,
		Username:  os.Getenv("SMTP_USERNAME"),
		Password:  os.Getenv("SMTP_PASSWORD"),
		FromEmail: os.Getenv("SMTP_FROM_EMAIL"),
		FromName:  os.Getenv("SMTP_FROM_NAME"),
		Enabled:   isEnabled,
	}
}

// EmailService handles email sending
type EmailService struct {
	config *EmailConfig
}

// NewEmailService creates a new email service instance
func NewEmailService() *EmailService {
	return &EmailService{
		config: GetEmailConfig(),
	}
}

// getTemplateFileName returns the template file name based on OTP type
func getTemplateFileName(otpType string) string {
	templateMap := map[string]string{
		"email-verification": "email_verification.html",
		"password-reset":     "password_reset.html",
		"2fa":                "two_factor_auth.html",
		"phone-verification": "phone_verification.html",
	}

	if filename, ok := templateMap[otpType]; ok {
		return filename
	}
	// Default to email verification if type not found
	return "email_verification.html"
}

// getEmailSubject returns the email subject based on OTP type
func getEmailSubject(otpType string) string {
	subjectMap := map[string]string{
		"email-verification": "Email Verification Code - Bixor Engine",
		"password-reset":     "Password Reset Code - Bixor Engine",
		"2fa":                "Two-Factor Authentication Code - Bixor Engine",
		"phone-verification": "Phone Verification Code - Bixor Engine",
	}

	if subject, ok := subjectMap[otpType]; ok {
		return subject
	}
	return "Verification Code - Bixor Engine"
}

// SendOTPEmail sends an OTP verification code email
func (es *EmailService) SendOTPEmail(otpType, toEmail, toName, otpCode string) error {
	if !es.config.Enabled {
		return fmt.Errorf("SMTP is not enabled. Set SMTP_ENABLED=true in .env file")
	}

	// Validate configuration
	if es.config.Host == "" || es.config.Username == "" || es.config.Password == "" {
		return fmt.Errorf("SMTP configuration is incomplete. Please check your .env file")
	}

	// Load template from embedded files
	templateFile := getTemplateFileName(otpType)
	// Use forward slashes for embedded filesystem (works on all platforms)
	templatePath := fmt.Sprintf("emails/templates/%s", templateFile)

	templateContent, err := emailTemplates.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load email template %s: %w", templateFile, err)
	}

	// Parse template
	tmpl, err := template.New("email").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	// Prepare data
	data := struct {
		FromEmail string
		FromName  string
		ToEmail   string
		ToName    string
		OTPCode   string
	}{
		FromEmail: es.config.FromEmail,
		FromName:  es.config.FromName,
		ToEmail:   toEmail,
		ToName:    toName,
		OTPCode:   otpCode,
	}

	// Execute template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	// Setup authentication
	auth := smtp.PlainAuth("", es.config.Username, es.config.Password, es.config.Host)

	// Email headers
	headers := fmt.Sprintf("From: %s <%s>\r\n", es.config.FromName, es.config.FromEmail)
	headers += fmt.Sprintf("To: %s <%s>\r\n", toName, toEmail)
	headers += fmt.Sprintf("Subject: %s\r\n", getEmailSubject(otpType))
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/html; charset=UTF-8\r\n"
	headers += "\r\n"

	// Combine headers and body
	message := []byte(headers + body.String())

	// Send email
	addr := fmt.Sprintf("%s:%d", es.config.Host, es.config.Port)
	err = smtp.SendMail(addr, auth, es.config.FromEmail, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// IsEnabled checks if email service is enabled
func (es *EmailService) IsEnabled() bool {
	return es.config.Enabled
}
