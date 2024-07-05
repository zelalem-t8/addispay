package addispay

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// AddisPay represents the AddisPay SDK
type AddisPay struct {
	PublicKey   string
	PrivateKey  string
	Auth        string
	CheckoutURL string
}

// New creates a new AddisPay instance
func New(publicKey, privateKey, auth string) *AddisPay {
	return &AddisPay{
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
		Auth:        auth,
		CheckoutURL: "https://uat-checkoutapi.addispay.et/api/v1/encrypted/receive-data/",
	}
}

func (ap *AddisPay) parsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	parsedKey, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return nil, err
	}
	return parsedKey.(*rsa.PublicKey), nil
}

func (ap *AddisPay) parsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	parsedKey, err := x509.ParsePKCS1PrivateKey(decodedPrivateKey)
	if err != nil {
		return nil, err
	}
	return parsedKey, nil
}

func (ap *AddisPay) encryptData(data string) string {
	rsaKey, err := ap.parsePublicKey(ap.PublicKey)
	if err != nil {
		return ""
	}
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, rsaKey, []byte(data))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes)
}

func (ap *AddisPay) decryptData(encryptedData string) (string, error) {
	rsaKey, err := ap.parsePrivateKey(ap.PrivateKey)
	if err != nil {
		return "", err
	}
	decodedEncryptedData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, rsaKey, decodedEncryptedData)
	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}

func (ap *AddisPay) sendRequest(totalAmount, txRef, currency, firstName, email, phoneNumber, lastName, sessionExpirationMinute, nonce, notifyURL, returnURL, message string) (*http.Response, error) {
	data := map[string]interface{}{
		"data": map[string]string{
			"total_amount":    ap.encryptData(totalAmount),
			"tx_ref":          ap.encryptData(txRef),
			"currency":        ap.encryptData(currency),
			"first_name":      ap.encryptData(firstName),
			"email":           ap.encryptData(email),
			"phone_number":    ap.encryptData(phoneNumber),
			"last_name":       ap.encryptData(lastName),
			"session_expired": ap.encryptData(sessionExpirationMinute),
			"nonce":           ap.encryptData(nonce),
			"order_detail":    `{"items": "rfid", "description": "I am testing this"}`,
			"notify_url":      ap.encryptData(notifyURL),
			"success_url":     ap.encryptData(returnURL),
			"cancel_url":      ap.encryptData(returnURL),
			"error_url":       ap.encryptData(returnURL),
		},
		"message": ap.encryptData(message),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ap.CheckoutURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Auth", ap.Auth)

	client := &http.Client{}
	return client.Do(req)
}
