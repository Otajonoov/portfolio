package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SMSClient struct {
	apiURL   string
	email    string
	password string
	headers  map[string]string
}

// NewSMSClient creates a new SMSClient instance.
func NewSMSClient(apiURL, email, password string) *SMSClient {
	return &SMSClient{
		apiURL:   apiURL,
		email:    email,
		password: password,
		headers:  make(map[string]string),
	}
}

// request sends an HTTP request and returns the response body as a byte slice.
func (c *SMSClient) request(method, apiPath string, data map[string]interface{}) (map[string]interface{}, error) {
	// Build the complete URL
	url := c.apiURL + apiPath

	// Set request headers
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// Check if an authorization token is present in the headers
	if authHeader, ok := c.headers["Authorization"]; ok {
		headers["Authorization"] = authHeader
	}

	// Serialize data to JSON
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// Set request headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}

	// Unmarshal JSON response
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return response, nil

}

// authUser performs user authentication and returns a token.
func (c *SMSClient) authUser() (string, error) {
	authdata := map[string]interface{}{
		"email":    c.email,
		"password": c.password,
	}

	// Send an authentication request
	authResp, err := c.request("POST", "auth/login", authdata)
	if err != nil {
		return "", err
	}

	token, ok := authResp["data"].(map[string]interface{})["token"].(string)
	if !ok {
		return "", fmt.Errorf("Authentication failed. Token not found in response.")
	}

	return token, nil
}

// sendSMS sends an SMS message.
func (c *SMSClient) sendSMS(phoneNumber, message string) (map[string]interface{}, error) {
	token, err := c.authUser()
	if err != nil {
		return nil, err
	}

	c.headers["Authorization"] = fmt.Sprintf("Bearer %s", token)

	data := map[string]interface{}{
		"from":         4546,
		"mobile_phone": phoneNumber,
		"message":      "stiv.uz uchun tasdiqlash kodi:\n" + message,
	}

	apiPath := "message/sms/send"

	return c.request("POST", apiPath, data)
}

func SMS(phone_number, message string) (map[string]string, error) {
	// Example usage of the SMSClient
	apiURL := "https://notify.eskiz.uz/api/"
	email := "otajonoov@gmail.com"
	password := "qgemXcWYIliNPQyvsfXI5ZBfVNPdnFN1madIP9gU"

	client := NewSMSClient(apiURL, email, password)

	data := map[string]string{
		phone_number: message,
	}

	// Authenticate to obtain a token
	token, err := client.authUser()
	if err != nil {
		fmt.Println("Authentication error:", err)
		return data, err
	}

	client.headers["Authorization"] = fmt.Sprintf("Bearer %s", token)

	smsResp, err := client.sendSMS(phone_number, message)
	if err != nil {
		fmt.Println("Send SMS error: ", err)
		fmt.Println("Send SMS error to: ", smsResp)
		return data, err
	}
	return data, nil

}
