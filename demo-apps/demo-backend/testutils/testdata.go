package testutils

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

const KeycloakDefaultRealm = "master"

func SetupTestRealm(ctx context.Context, realmName, keycloakAdminUsername, keycloakAdminPassword, keycloakBasePath string) error {
	gocloakClient := gocloak.NewClient(keycloakBasePath)
	token, err := gocloakClient.LoginAdmin(
		ctx,
		keycloakAdminUsername,
		keycloakAdminPassword,
		KeycloakDefaultRealm,
	)
	if err != nil {
		return err
	}

	realmR := gocloak.RealmRepresentation{
		Realm:       gocloak.StringP(realmName),
		DisplayName: gocloak.StringP(realmName),
		Enabled:     gocloak.BoolP(true),
	}
	_, err = gocloakClient.CreateRealm(ctx, token.AccessToken, realmR)
	if err != nil {
		return err
	}
	clientID, clientRoles, err := createClientWithClientRoles(ctx, gocloakClient, token, realmName)
	if err != nil {
		return err
	}

	group := gocloak.Group{
		Name: gocloak.StringP("demo-group"),
	}
	groupID, err := gocloakClient.CreateGroup(ctx, token.AccessToken, realmName, group)
	if err != nil {
		return err
	}
	err = gocloakClient.AddClientRolesToGroup(ctx, token.AccessToken, realmName, clientID, groupID, clientRoles)
	if err != nil {
		return err
	}

	userID, err := createTestUser(ctx, gocloakClient, token, realmName, "test-user-1")
	if err != nil {
		return err
	}

	serviceAccount, err := gocloakClient.GetClientServiceAccount(ctx, token.AccessToken, realmName, userID)
	if err != nil {
		return err
	}
	if serviceAccount == nil {
		return fmt.Errorf("service account with ID %s not found", userID)
	}

	err = gocloakClient.AddUserToGroup(ctx, token.AccessToken, realmName, *serviceAccount.ID, groupID)
	if err != nil {
		return err
	}

	return nil
}

func createTestUser(ctx context.Context, gocloakClient *gocloak.GoCloak, token *gocloak.JWT, realmName string, username string) (string, error) {
	clientConfig := gocloak.Client{
		ClientID:                     gocloak.StringP(username),
		Name:                         gocloak.StringP(username),
		Enabled:                      gocloak.BoolP(true),
		Protocol:                     gocloak.StringP("openid-connect"),
		ClientAuthenticatorType:      gocloak.StringP("client-secret"),
		ServiceAccountsEnabled:       gocloak.BoolP(true),
		PublicClient:                 gocloak.BoolP(false),
		StandardFlowEnabled:          gocloak.BoolP(false),
		DirectAccessGrantsEnabled:    gocloak.BoolP(false),
		AuthorizationServicesEnabled: gocloak.BoolP(false),
	}
	clientID, err := gocloakClient.CreateClient(ctx, token.AccessToken, realmName, clientConfig)
	if err != nil {
		return "", err
	}

	client, err := gocloakClient.GetClient(ctx, token.AccessToken, realmName, clientID)
	if err != nil {
		return "", err
	}

	client.Secret = gocloak.StringP(username + "-secret")
	err = gocloakClient.UpdateClient(ctx, token.AccessToken, realmName, *client)
	if err != nil {
		return "", err
	}

	return clientID, nil
}

func createClientWithClientRoles(ctx context.Context, gocloakClient *gocloak.GoCloak, token *gocloak.JWT, realmName string) (string, []gocloak.Role, error) {
	rediectURIs := []string{"http://localhost:*"}
	clientConfig := gocloak.Client{
		ClientID:                  gocloak.StringP("demo-app"),
		Name:                      gocloak.StringP("demo-app"),
		Enabled:                   gocloak.BoolP(true),
		StandardFlowEnabled:       gocloak.BoolP(true),
		ImplicitFlowEnabled:       gocloak.BoolP(false),
		DirectAccessGrantsEnabled: gocloak.BoolP(true),
		ServiceAccountsEnabled:    gocloak.BoolP(true),
		WebOrigins:                &rediectURIs,
		PublicClient:              gocloak.BoolP(false),
		Attributes: &map[string]string{
			"policy_enforcement_mode": "ENFORCING",
		},
	}
	clientID, err := gocloakClient.CreateClient(ctx, token.AccessToken, realmName, clientConfig)
	if err != nil {
		return "", nil, err
	}

	client, err := gocloakClient.GetClient(ctx, token.AccessToken, realmName, clientID)
	if err != nil {
		return "", nil, err
	}

	newSecret := "demo-app-secret"
	client.Secret = &newSecret

	err = gocloakClient.UpdateClient(ctx, token.AccessToken, realmName, *client)
	if err != nil {
		return "", nil, err
	}

	clientRoleAConfig := gocloak.Role{
		Name:       gocloak.StringP("demo.role.a"),
		ClientRole: gocloak.BoolP(true),
	}
	_, err = gocloakClient.CreateClientRole(ctx, token.AccessToken, realmName, *client.ID, clientRoleAConfig)
	if err != nil {
		return "", nil, err
	}

	roleA, err := gocloakClient.GetClientRole(ctx, token.AccessToken, realmName, *client.ID, "demo.role.a")
	if err != nil {
		return "", nil, err
	}
	clientRoles := []gocloak.Role{
		*roleA,
	}

	scopes, err := gocloakClient.GetClientScopes(ctx, token.AccessToken, realmName)
	if err != nil {
		return "", nil, err
	}
	var scopeRoles *gocloak.ClientScope
	for i, s := range scopes {
		fmt.Println(i, s)
		if *s.Name == "roles" {
			scopeRoles = s
		}
	}

	var pmClientRoles *gocloak.ProtocolMappers
	pms, err := gocloakClient.GetClientScopeProtocolMappers(ctx, token.AccessToken, realmName, *scopeRoles.ID)
	if err != nil {
		return "", nil, err
	}
	for i, pm := range pms {
		fmt.Println(i, pm)
		if *pm.Name == "client roles" {
			pm.ProtocolMappersConfig = &gocloak.ProtocolMappersConfig{
				ClaimName:          gocloak.StringP("roles"),
				UserAttribute:      gocloak.StringP("roles"),
				UserinfoTokenClaim: gocloak.StringP("true"),
				Multivalued:        gocloak.StringP("true"),
				IDTokenClaim:       gocloak.StringP("true"),
				JSONTypeLabel:      gocloak.StringP("String"),
				AccessTokenClaim:   gocloak.StringP("true"),
			}
			pmClientRoles = pm
		}
	}

	err = gocloakClient.UpdateClientScopeProtocolMapper(ctx, token.AccessToken, realmName, *scopeRoles.ID, *pmClientRoles)
	if err != nil {
		return "", nil, err
	}

	return *client.ID, clientRoles, nil
}
