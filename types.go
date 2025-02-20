package main

// EmailData represents the data needed for sending emails
type EmailData struct {
	CustomerName  string
	CustomerEmail string
	CourseName    string
	CompanyName   string
	SupportEmail  string
	DomainURL     string
	SenderEmail   string
}

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// Email represents an email to be sent
type Email struct {
	To          string
	From        string
	Subject     string
	HTMLContent string
	TextContent string
}
