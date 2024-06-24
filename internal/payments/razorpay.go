package payments

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	// "errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

// PaymentRequest represents the expected input for the payment initiation
type PaymentRequest struct {
	Amount   int    `json:"amount" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Receipt  string `json:"receipt" binding:"required"`
}

type PaymentVerificationRequest struct {
	Signature string `json:"signature" binding:"required"`
	OrderID   string `json:"order_id" binding:"required"`
	PaymentID string `json:"payment_id" binding:"required"`
}

// executerazorpay handles the payment initiation via Razorpay
func executerazorpay(c *gin.Context) {
	var req PaymentRequest

	// Bind JSON input to PaymentRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := razorpay.NewClient("rzp_test_CQh1HvVjhsOvaR", "cJM6USA5rdAEmLG4Dzq5qIiI")

	data := map[string]interface{}{
		"amount":   req.Amount * 100, // Amount in paisa
		"currency": req.Currency,
		"receipt":  req.Receipt,
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		log.Println("Error creating Razorpay order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment not initiated"})
		return
	}

	razorId, ok := body["id"].(string)
	if !ok {
		log.Println("Error extracting Razorpay order ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment not initiated"})
		return
	}

	log.Println("Razorpay Order ID:", razorId)
	c.JSON(http.StatusOK, gin.H{"razorpay_order_id": razorId})
}

   func RazorPaymentVerification(c *gin.Context) {
	var req PaymentVerificationRequest

	// Bind JSON input to PaymentVerificationRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	secret := "cJM6USA5rdAEmLG4Dzq5qIiI"
	data := req.OrderID + "|" + req.PaymentID

	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(data))
	if err != nil {
		log.Println("Error generating HMAC:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	sha := hex.EncodeToString(h.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(sha), []byte(req.Signature)) != 1 {
		log.Println("Failed payment verification")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "payment verification failed"})
		return
	}

	log.Println("Payment verification successful")
	c.JSON(http.StatusOK, gin.H{"message": "payment verification successful"})
}


func RegisterPaymentRoutes(rg *gin.RouterGroup) {
	paymentroute := rg.Group("/payments")
	paymentroute.POST("/initiate", executerazorpay)
	paymentroute.POST("/verify", RazorPaymentVerification)
}


