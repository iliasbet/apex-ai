package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/product"
)

// EmailData represents the data needed for sending emails
type EmailData struct {
	CustomerName  string
	CustomerEmail string
	CourseName    string
	CompanyName   string
	SupportEmail  string
}

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Set your secret key
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Fatal("STRIPE_SECRET_KEY is required")
	}
	stripe.Key = stripeKey
}

// sendWelcomeEmail sends a welcome email to the customer using SMTP
func sendWelcomeEmail(data EmailData) error {
	config := EmailConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SENDER_EMAIL"),
	}

	// Create HTML email template
	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
				.container { max-width: 600px; margin: 0 auto; padding: 20px; }
				.header { text-align: center; padding: 20px 0; }
				.content { background: #f9f9f9; padding: 20px; border-radius: 5px; }
				.button { display: inline-block; padding: 12px 24px; background: #0066FF; color: white; text-decoration: none; border-radius: 5px; }
				.footer { text-align: center; padding: 20px 0; font-size: 0.9em; color: #666; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Welcome to %s!</h1>
				</div>
				<div class="content">
					<p>Dear %s,</p>
					<p>Thank you for enrolling in %s! We're excited to have you join us on this transformative journey into AI implementation and strategy.</p>
					<p>Here's what happens next:</p>
					<ol>
						<li>You'll receive your login credentials within the next 10 minutes</li>
						<li>Access to all course materials will be immediately available</li>
						<li>You can start your learning journey right away</li>
					</ol>
					<p style="text-align: center; margin: 30px 0;">
						<a href="%s/login" class="button">Access Your Course</a>
					</p>
					<p>If you need any assistance, our support team is here to help at <a href="mailto:%s">%s</a>.</p>
				</div>
				<div class="footer">
					<p>© %s. All rights reserved.</p>
					<p>This email was sent to %s</p>
				</div>
			</div>
		</body>
		</html>
	`, data.CourseName, data.CustomerName, data.CourseName, os.Getenv("DOMAIN_URL"),
		data.SupportEmail, data.SupportEmail, data.CompanyName, data.CustomerEmail)

	// Set email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", os.Getenv("SENDER_NAME"), config.From)
	headers["To"] = data.CustomerEmail
	headers["Subject"] = fmt.Sprintf("Welcome to %s!", data.CourseName)
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build email message
	var message bytes.Buffer
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(htmlContent)

	// Connect to SMTP server
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	// Send email
	return smtp.SendMail(
		addr,
		auth,
		config.From,
		[]string{data.CustomerEmail},
		message.Bytes(),
	)
}

// createOrGetProduct ensures our product exists in Stripe
func createOrGetProduct() (*stripe.Product, error) {
	// Try to find existing product
	listParams := &stripe.ProductListParams{}
	listParams.Filters.AddFilter("metadata[product_id]", "", "apex_ai_course")

	products := product.List(listParams)
	for products.Next() {
		return products.Product(), nil
	}

	// Create new product if not found
	productParams := &stripe.ProductParams{
		Name:        stripe.String(os.Getenv("COURSE_NAME")),
		Description: stripe.String("Complete AI transformation course for business leaders"),
		Images:      []*string{stripe.String("https://your-domain.com/course-image.jpg")},
		URL:         stripe.String("https://your-domain.com/course"),
		DefaultPriceData: &stripe.ProductDefaultPriceDataParams{
			UnitAmount: stripe.Int64(299900), // $2,999.00
			Currency:   stripe.String(string(stripe.CurrencyUSD)),
		},
	}
	productParams.AddMetadata("product_id", "apex_ai_course")

	return product.New(productParams)
}

// PaymentHandler creates a Stripe checkout session and redirects to it
func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	// Get or create the product
	prod, err := createOrGetProduct()
	if err != nil {
		log.Printf("Error creating/getting product: %v", err)
		http.Error(w, "Error setting up payment", http.StatusInternalServerError)
		return
	}

	// Create checkout session with enhanced customization
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(fmt.Sprintf("%s/payment-success?session_id={CHECKOUT_SESSION_ID}", os.Getenv("DOMAIN_URL"))),
		CancelURL:  stripe.String(fmt.Sprintf("%s/payment", os.Getenv("DOMAIN_URL"))),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(prod.DefaultPrice.ID),
				Quantity: stripe.Int64(1),
			},
		},
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
			"link", // Add support for Link by Stripe
		}),
		AllowPromotionCodes:      stripe.Bool(true),
		BillingAddressCollection: stripe.String(string(stripe.CheckoutSessionBillingAddressCollectionRequired)),
		CustomerCreation:         stripe.String(string(stripe.CheckoutSessionCustomerCreationAlways)),
		PhoneNumberCollection: &stripe.CheckoutSessionPhoneNumberCollectionParams{
			Enabled: stripe.Bool(true),
		},
		CustomFields: []*stripe.CheckoutSessionCustomFieldParams{
			{
				Key: stripe.String("company_name"),
				Label: &stripe.CheckoutSessionCustomFieldLabelParams{
					Type:   stripe.String("custom"),
					Custom: stripe.String("Company Name"),
				},
				Type:     stripe.String(string(stripe.CheckoutSessionCustomFieldTypeText)),
				Optional: stripe.Bool(false),
			},
			{
				Key: stripe.String("job_title"),
				Label: &stripe.CheckoutSessionCustomFieldLabelParams{
					Type:   stripe.String("custom"),
					Custom: stripe.String("Job Title"),
				},
				Type:     stripe.String(string(stripe.CheckoutSessionCustomFieldTypeText)),
				Optional: stripe.Bool(false),
			},
		},
		CustomText: &stripe.CheckoutSessionCustomTextParams{
			ShippingAddress: &stripe.CheckoutSessionCustomTextShippingAddressParams{
				Message: stripe.String("Please provide your business address for billing purposes."),
			},
			Submit: &stripe.CheckoutSessionCustomTextSubmitParams{
				Message: stripe.String("By completing this purchase, you agree to our Terms of Service and Privacy Policy."),
			},
		},
	}

	session, err := session.New(params)
	if err != nil {
		log.Printf("Error creating checkout session: %v", err)
		http.Error(w, "Error setting up payment", http.StatusInternalServerError)
		return
	}

	// Redirect to Stripe's checkout page
	http.Redirect(w, r, session.URL, http.StatusSeeOther)
}

// PaymentSuccessHandler handles the success page after payment
func PaymentSuccessHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Invalid session", http.StatusBadRequest)
		return
	}

	// Verify the session
	checkoutSession, err := session.Get(sessionID, nil)
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Error verifying payment", http.StatusInternalServerError)
		return
	}

	// Send welcome email
	emailData := EmailData{
		CustomerName:  checkoutSession.CustomerDetails.Name,
		CustomerEmail: checkoutSession.CustomerEmail,
		CourseName:    os.Getenv("COURSE_NAME"),
		CompanyName:   os.Getenv("COMPANY_NAME"),
		SupportEmail:  os.Getenv("SUPPORT_EMAIL"),
	}

	if err := sendWelcomeEmail(emailData); err != nil {
		log.Printf("Error sending welcome email: %v", err)
		// Continue anyway, as the payment was successful
	}

	// Show success page
	w.Write([]byte(fmt.Sprintf(`
		<html>
			<head>
				<title>Payment Successful</title>
				<link href="https://fonts.googleapis.com/css2?family=Lexend+Deca:wght@400;500;600;700&display=swap" rel="stylesheet">
				<script src="https://cdn.tailwindcss.com"></script>
			</head>
			<body class="bg-black text-white min-h-screen flex items-center justify-center font-['Lexend_Deca']">
				<div class="text-center p-8">
					<h1 class="text-4xl font-bold mb-4">Payment Successful!</h1>
					<p class="text-xl text-blue-200 mb-8">Thank you for enrolling in %s!</p>
					<p class="text-lg text-blue-200/90 mb-4">We've sent your course access details to %s</p>
					<p class="text-sm text-blue-200/70">If you don't receive the email within 10 minutes, please check your spam folder or contact <a href="mailto:%s" class="text-blue-400 hover:text-blue-300">%s</a></p>
				</div>
			</body>
		</html>
	`, os.Getenv("COURSE_NAME"), checkoutSession.CustomerEmail, os.Getenv("SUPPORT_EMAIL"), os.Getenv("SUPPORT_EMAIL"))))
}
