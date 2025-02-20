package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/product"
)

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

	// Initialize email service
	emailService := NewEmailService()

	// Send welcome email
	emailData := EmailData{
		CustomerName:  checkoutSession.CustomerDetails.Name,
		CustomerEmail: checkoutSession.CustomerEmail,
		CourseName:    os.Getenv("COURSE_NAME"),
		CompanyName:   os.Getenv("COMPANY_NAME"),
		SupportEmail:  os.Getenv("SUPPORT_EMAIL"),
	}

	if err := emailService.SendWelcomeEmail(emailData); err != nil {
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
