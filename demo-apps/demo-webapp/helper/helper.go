package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func FormatJSON(jsonData []byte) ([]byte, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, jsonData, "", "  ") // Two spaces for indentation
	if err != nil {
		return nil, err
	}
	return prettyJSON.Bytes(), nil
}

// DecodeToken takes a JWT token string and returns a map of its claims
func DecodeToken(tokenString string) (map[string]interface{}, error) {
	// Split the token into its three parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode the payload (second part) from base64
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding payload: %v", err)
	}

	// Parse the JSON payload into a map
	var claims map[string]interface{}
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return nil, fmt.Errorf("error parsing claims: %v", err)
	}

	// Return the parsed claims
	return claims, nil
}

// DecodeIDToken decode and parse the ID token
func DecodeIDToken(idToken string) (map[string]interface{}, error) {
	// Split the token into Header, Payload, and Signature
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode the Payload (the second part of the token)
	payload, err := decodeBase64URL(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	// Parse the Payload as JSON into a map
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse payload JSON: %v", err)
	}

	return claims, nil
}

// decodeBase64URL decode a Base64 URL-safe encoded string (with padding)
func decodeBase64URL(input string) ([]byte, error) {
	// Replace URL-safe characters with standard Base64 characters
	decodedInput := strings.Replace(input, "-", "+", -1)
	decodedInput = strings.Replace(decodedInput, "_", "/", -1)

	// Add missing padding if necessary
	switch len(decodedInput) % 4 {
	case 2:
		decodedInput += "=="
	case 3:
		decodedInput += "="
	}

	// Decode the Base64 string
	return base64.StdEncoding.DecodeString(decodedInput)
}
