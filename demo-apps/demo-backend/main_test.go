package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"slices"
	"testing"

	"github.com/malfter/keycloak/demo-applications/demo-backend/testutils"
	"github.com/malfter/keycloak/demo-applications/demo-backend/utils"
	keycloak "github.com/stillya/testcontainers-keycloak"
)

var (
	keycloakImage         = "keycloak/keycloak:26.1.3"
	keycloakAdminUsername = "admin"
	keycloakAdminPassword = "admin"
	keycloakTestRealm     = "keycloak-demo"
)

func TestAuthMiddlewareHandler(t *testing.T) {
	// Setup
	ctx := context.Background()
	authServerURL, issuerURL, err := setupAuthServer(ctx, t)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Initialize verifier and handler
	clientID := "demo-app"
	initVerifier(issuerURL, clientID)
	handler := authMiddleware(http.HandlerFunc(serviceHandler))

	// Get fake token
	fakeToken, err := getFakeToken(authServerURL, t)
	if err != nil {
		t.Fatalf("Failed to get fake token: %v", err)
	}

	// Prepare request
	req, w := prepareRequest(fakeToken)

	// Execute request
	handler.ServeHTTP(w, req)

	// Verify response
	verifyResponse(t, w)

	// Verify claims and roles
	verifyClaims(t, w)
}

func setupAuthServer(ctx context.Context, t *testing.T) (string, string, error) {
	authServerURL, err := keycloakContainer.GetAuthServerURL(ctx)
	if err != nil {
		return "", "", fmt.Errorf("GetAuthServerURL() error: %v", err)
	}

	baseURL, err := url.Parse(authServerURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid base URL error: %v", err)
	}

	issuerURL := baseURL.JoinPath("realms", keycloakTestRealm).String()
	return authServerURL, issuerURL, nil
}

func getFakeToken(authServerURL string, t *testing.T) (string, error) {
	return testutils.GetToken(
		authServerURL,
		keycloakTestRealm,
		"test-user-1",
		"test-user-1-secret",
	)
}

func prepareRequest(token string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req, httptest.NewRecorder()
}

func verifyResponse(t *testing.T, w *httptest.ResponseRecorder) {
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 but got %d", w.Code)
	}

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	if err != nil {
		t.Errorf("unmarshal body failed with error %v", err)
	}

	expectedMessage := "Authenticated successfully"
	if message := body["message"]; message != expectedMessage {
		t.Errorf("expected response %s but got %s", expectedMessage, message)
	}
}

func verifyClaims(t *testing.T, w *httptest.ResponseRecorder) {
	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	if err != nil {
		t.Errorf("unmarshal body failed with error %v", err)
	}

	claims := body["claims"].(map[string]interface{})
	roles, err := utils.ExtractRoles(claims)
	if err != nil {
		t.Errorf("extract roles from claims failed error = %v", err)
	}

	expectedRole := "demo.role.a"
	if !slices.Contains(roles, expectedRole) {
		t.Errorf("roles does not contain expected role %s (roles: %s)", expectedRole, roles)
	}
}

var keycloakContainer *keycloak.KeycloakContainer

func setup() {
	var err error
	ctx := context.Background()
	keycloakContainer, err = RunContainer(ctx)
	if err != nil {
		panic(err)
	}
	keycloakServerURL, err := keycloakContainer.GetAuthServerURL(ctx)
	if err != nil {
		panic(err)
	}
	err = testutils.SetupTestRealm(
		ctx,
		keycloakTestRealm,
		keycloakAdminUsername,
		keycloakAdminPassword,
		keycloakServerURL,
	)
	if err != nil {
		panic(err)
	}
}

func shutDown() {
	ctx := context.Background()
	err := keycloakContainer.Terminate(ctx)
	if err != nil {
		panic(err)
	}
}

func RunContainer(ctx context.Context) (*keycloak.KeycloakContainer, error) {
	return keycloak.Run(ctx,
		keycloakImage,
		keycloak.WithContextPath("/auth"),
		keycloak.WithAdminUsername(keycloakAdminUsername),
		keycloak.WithAdminPassword(keycloakAdminPassword),
	)
}

func Test_Example_PrintServerURL(t *testing.T) {
	ctx := context.Background()

	authServerURL, err := keycloakContainer.GetAuthServerURL(ctx)
	if err != nil {
		t.Errorf("GetAuthServerURL() error = %v", err)
		return
	}

	fmt.Println(authServerURL)
	// Output:
	// http://localhost:32768/auth
}

func TestMain(m *testing.M) {
	defer func() {
		if r := recover(); r != nil {
			shutDown()
			fmt.Printf("Panic: %s", r)
		}
	}()
	setup()
	code := m.Run()
	shutDown()
	os.Exit(code)
}
