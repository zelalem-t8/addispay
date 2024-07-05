package main

import (
	"fmt"

	"github.com/zelalem-t8/addispay" // Import your SDK package
)

func main() {
	// Initialize AddisPay instance
	publicKey := "your-public-key"
	privateKey := "your-private-key"
	auth := "your-auth-token"
	ap := addispay.New(publicKey, privateKey, auth)

	// Example call to SendRequest
	resp, err := ap.SendRequest("100.0", "tx123", "USD", "John", "john@example.com", "123456789", "Doe", "30", "nonce123", "https://example.com/notify", "https://example.com/return", "Test message")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()
	fmt.Println("Response Status:", resp.Status)
	// Handle response data as needed
}
