package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/malfter/keycloak/demo-applications/demo-webapp/helper"
	"golang.org/x/oauth2"
)

var oauth2Config oauth2.Config
var state = "random-state" // In production use a secure random state
var backendServiceURL = "http://localhost:9082/"

// Initialize OAuth2 config for the Go web app to communicate with Keycloak
func init() {
	oauth2Config = oauth2.Config{
		ClientID:     "demo-webapp", // Your Keycloak client ID
		ClientSecret: "demo-webapp", // Your Keycloak client secret
		RedirectURL:  "http://localhost:9081/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:9080/realms/keycloak-demo/protocol/openid-connect/auth",
			TokenURL: "http://localhost:9080/realms/keycloak-demo/protocol/openid-connect/token",
		},
		Scopes: []string{"openid", "profile", "email", "offline_access"},
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	html := `
	<!DOCTYPE html>
		<html>
			<head>
				<title>Demo Webapp</title>
			</head>
		<body>
			<h1>Demo Webapp</h1>
			<p>
				<a href="/login">Login</a>
			</p>
		</body>
	</html>
	`

	fmt.Fprintf(w, "%s", html)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a URL to Keycloak login
	authURL := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the callback from Keycloak
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code is missing", http.StatusBadRequest)
		return
	}

	// Exchange the code for a token
	token, err := oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "Access Token: %s\n\n", token.AccessToken)
	// claims, err := helper.DecodeToken(token.AccessToken)
	// if err != nil {
	// 	log.Fatalf("Error decoding Access Token: %v", err)
	// }
	// fmt.Fprint(w, "Decoded Access Token:\n")
	// jsonData, _ := json.MarshalIndent(claims, "", "  ")
	// fmt.Fprintf(w, "%s\n\n", jsonData)

	// Retrieve the id_token
	idToken := token.Extra("id_token").(string)
	if idToken != "" {
		// fmt.Fprintf(w, "ID Token (raw): %s\n\n", idToken)

		// claims, err := helper.DecodeIDToken(idToken)
		// if err != nil {
		// 	log.Fatalf("Error decoding ID Token: %v", err)
		// }
		// fmt.Fprint(w, "Decoded ID Token:\n")
		// jsonData, _ := json.MarshalIndent(claims, "", "  ")
		// fmt.Fprintf(w, "%s\n\n", jsonData)
	} else {
		http.Error(w, "ID Token is missing", http.StatusInternalServerError)
		return
	}

	callBackendService(w, idToken)
}

func callBackendService(w http.ResponseWriter, bearerToken string) {
	// The URL to which the GET request will be sent
	url := backendServiceURL

	// Your Bearer Token
	token := bearerToken

	// Create the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Authorization header with the Bearer Token
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request with an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close() // Ensure the response body is closed after reading

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyPretty, err := helper.FormatJSON(body)
	if err != nil {
		log.Fatal(err)
	}

	// Output the status code and the response body
	fmt.Fprintf(w, "Status Code: %s\n", resp.Status)
	fmt.Fprintf(w, "Response Body: %s\n", "")
	fmt.Fprintf(w, "%s\n\n", string(bodyPretty))
}

func main() {
	// Define routes
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)

	// Start the web app server
	log.Println("Web app is running on http://localhost:9081")
	log.Fatal(http.ListenAndServe(":9081", nil))
}
