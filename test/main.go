package main

import (
	"fmt"

	"github.com/zelalem-t8/addispay"
)

func main() {
	addisPay := addispay.New("your_base64_public_key", "your_base64_private_key", "your_auth_token")

	response, err := addisPay.sendRequest("1000", "12345", "ETB", "John", "john@example.com", "0912345678", "Doe", "30", "unique_nonce", "https://your.notify.url", "https://your.return.url", "Test transaction")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response:", response.Status)
	}
}
