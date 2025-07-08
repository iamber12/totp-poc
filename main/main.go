package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"net/http"
	"time"
)

var secret string

func main() {
	r := gin.Default()

	// Step 1: Register — generates secret & saves QR code image
	r.GET("/register", func(c *gin.Context) {
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "GoTOTPApp",
			AccountName: "user@example.com",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
			return
		}

		secret = key.Secret()

		// Save QR code directly to disk
		if err := qrcode.WriteFile(key.URL(), qrcode.Medium, 256, "totp-qr.png"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save QR code image"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     "Scan the QR code from 'totp-qr.png'",
			"secret":      secret,
			"otpauth_url": key.URL(),
			"qr_file":     "Saved to totp-qr.png",
		})
	})

	// Step 2: Verify OTP from user
	r.POST("/verify", func(c *gin.Context) {
		var req struct {
			OTP string `json:"otp"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if secret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Call /register first"})
			return
		}

		if valid := totp.Validate(req.OTP, secret); !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "❌ Invalid OTP"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "✅ OTP is valid"})
	})

	// Debug route to print the expected OTP at server side
	r.GET("/debug-otp", func(c *gin.Context) {
		if secret == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Call /register first"})
			return
		}
		code, _ := totp.GenerateCode(secret, time.Now())
		c.JSON(http.StatusOK, gin.H{"expected_otp": code})
	})

	r.Run(":8080")
}
