package testutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Global HTTP client with a timeout to prevent indefinite waits
var httpClient = &http.Client{Timeout: 10 * time.Second}

// GetToken retrieves an access token from Keycloak using client credentials.
func GetToken(keycloakBaseURL, realm, clientID, clientSecret string) (string, error) {
	// Parse the base URL to ensure it's valid
	baseURL, err := url.Parse(keycloakBaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	// Construct the full URL for the token endpoint
	tokenURL := baseURL.JoinPath("realms", realm, "protocol/openid-connect/token")

	// Prepare the form data for the request
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	// Create a new POST request with the form data
	req, err := http.NewRequest("POST", tokenURL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	// Set the content type to match the form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request and handle any errors
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	// Ensure the response body is closed after use
	defer resp.Body.Close()

	// Check if the response status is OK (200)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	// Decode the JSON response into a TokenResponse struct
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	// Return the access token from the response
	return tokenResponse.AccessToken, nil
}
