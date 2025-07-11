# Go TOTP Proof of Concept

This project demonstrates a simple Time-based One-Time Password (TOTP) authentication flow using Go, Gin, and QR code generation.

## Features

- Register endpoint to generate a TOTP secret and QR code
- Verify endpoint to validate user-provided OTPs
- Debug endpoint to view the current valid OTP (for testing)

## Endpoints

### 1. Register

- **GET** `/register`
- Generates a new TOTP secret and saves a QR code as `totp-qr.png`.
- Response includes the secret and otpauth URL.

### 2. Verify

- **POST** `/verify`
- Request body: `{ "otp": "<6-digit code>" }`
- Validates the provided OTP against the current secret.

### 3. Debug OTP

- **GET** `/debug-otp`
- Returns the current valid OTP for the secret (for testing only).

## Usage

1. Start the server:
   ```sh
   go run main/main.go