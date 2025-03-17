package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
)

type contextKey string

const (
	issuerURL            = "http://localhost:9080/realms/keycloak-demo"
	clientID             = "demo-webapp"
	claimsKey contextKey = "claims"
)

var verifier *oidc.IDTokenVerifier

func init() {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	verifier = provider.Verifier(&oidc.Config{ClientID: clientID})
}

// Middleware to check Authorization header
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" || len(tokenString) < 8 || !strings.HasPrefix(tokenString, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[7:]
		ctx := r.Context()
		token, err := verifier.Verify(ctx, tokenString)
		if err != nil {
			fmt.Println("error:", err)
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Extract user information from token
		var claims map[string]interface{}
		if err := token.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		// Store claims in the request context
		ctx = context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(claimsKey).(map[string]interface{})
	response := map[string]interface{}{
		"message": "Authenticated successfully",
		"claims":  claims,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/", authMiddleware(serviceHandler))

	fmt.Println("Server running on port 9082")
	log.Fatal(http.ListenAndServe(":9082", nil))
}
