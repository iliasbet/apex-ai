package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"time"
)

// EmailService handles email operations
type EmailService struct {
	config EmailConfig
}

// NewEmailService creates a new email service instance
func NewEmailService() *EmailService {
	return &EmailService{
		config: EmailConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SENDER_EMAIL"),
		},
	}
}

// SendWelcomeEmail sends a welcome email to the customer
func (s *EmailService) SendWelcomeEmail(data EmailData) error {
	// Add domain URL and sender email to the template data
	data.DomainURL = os.Getenv("DOMAIN_URL")
	data.SenderEmail = s.config.From

	// Create welcome email content
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        /* Base styles */
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
            background-color: #f9f9f9;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            padding: 20px 0;
            background-color: #0066FF;
            color: white;
            border-radius: 8px 8px 0 0;
        }
        .content {
            padding: 30px 20px;
            background: #ffffff;
        }
        .button {
            display: inline-block;
            padding: 12px 24px;
            background: #0066FF;
            color: white !important;
            text-decoration: none;
            border-radius: 5px;
            font-weight: bold;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            padding: 20px 0;
            font-size: 0.9em;
            color: #666;
            border-top: 1px solid #eee;
        }
        .next-steps {
            background: #f5f7ff;
            padding: 20px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .next-steps h3 {
            margin-top: 0;
            color: #0066FF;
        }
        .support-section {
            background: #fff8f0;
            padding: 15px;
            border-radius: 5px;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to {{.CourseName}}!</h1>
        </div>
        <div class="content">
            <p>Dear {{.CustomerName}},</p>
            
            <p>Thank you for enrolling in <strong>{{.CourseName}}</strong>! We're excited to have you join us on this transformative journey into AI implementation and strategy.</p>
            
            <div class="next-steps">
                <h3>🚀 Here's what happens next:</h3>
                <ol>
                    <li>Check your inbox for your login credentials (arriving within 10 minutes)</li>
                    <li>Access the complete course materials immediately after login</li>
                    <li>Join our community of business leaders and AI innovators</li>
                    <li>Start your learning journey at your own pace</li>
                </ol>
            </div>

            <p style="text-align: center;">
                <a href="{{.DomainURL}}/login" class="button">Access Your Course</a>
            </p>

            <div class="support-section">
                <p><strong>Need Help?</strong></p>
                <p>Our support team is here for you at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a></p>
                <p>We typically respond within 2 hours during business hours.</p>
            </div>
        </div>
        
        <div class="footer">
            <p>© {{.CompanyName}}. All rights reserved.</p>
            <p>This email was sent to {{.CustomerEmail}}</p>
            <p><small>Please add {{.SenderEmail}} to your contacts to ensure you receive our communications.</small></p>
        </div>
    </div>
</body>
</html>`

	// Parse the template
	tmpl, err := template.New("welcome").Parse(htmlContent)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Execute the template with the data
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	// Create email
	email := &Email{
		To:          data.CustomerEmail,
		From:        fmt.Sprintf("%s <%s>", os.Getenv("SENDER_NAME"), s.config.From),
		Subject:     fmt.Sprintf("Welcome to %s!", data.CourseName),
		HTMLContent: body.String(),
	}

	// Send with retry
	return s.sendEmailWithRetry(email, 3)
}

// sendEmailWithRetry attempts to send an email with retries
func (s *EmailService) sendEmailWithRetry(email *Email, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := s.sendEmail(email); err != nil {
			lastErr = err
			// Wait before retrying (exponential backoff)
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to send email after %d attempts: %v", maxRetries, lastErr)
}

// sendEmail sends a single email
func (s *EmailService) sendEmail(email *Email) error {
	// Set email headers
	headers := make(map[string]string)
	headers["From"] = email.From
	headers["To"] = email.To
	headers["Subject"] = email.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build email message
	var message bytes.Buffer
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(email.HTMLContent)

	// Connect to SMTP server
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	// Send email
	return smtp.SendMail(
		addr,
		auth,
		s.config.From,
		[]string{email.To},
		message.Bytes(),
	)
}
