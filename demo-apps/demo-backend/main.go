package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
)

const (
	issuerURL = "http://localhost:9080/realms/keycloak-demo"
	clientID  = "demo-webapp"
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
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || len(tokenString) < 8 || tokenString[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenString = tokenString[7:]
		ctx := context.Background()
		token, err := verifier.Verify(ctx, tokenString)
		if err != nil {
			fmt.Println("error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Extract user information from token
		var claims map[string]interface{}
		if err := token.Claims(&claims); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func serviceHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	c.JSON(http.StatusOK, gin.H{
		"message": "Authenticated successfully",
		"claims":  claims,
	})
}

func main() {
	r := gin.Default()
	r.GET("/", authMiddleware(), serviceHandler)

	fmt.Println("Server running on port 9082")
	r.Run(":9082")
}
