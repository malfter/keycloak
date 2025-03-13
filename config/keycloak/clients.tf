resource "keycloak_openid_client" "webapp" {
  realm_id                     = keycloak_realm.this.id
  client_id                    = var.webapp_client_id
  client_secret                = var.webapp_client_id
  name                         = "Web Application"
  enabled                      = true
  access_type                  = "CONFIDENTIAL"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = true
  service_accounts_enabled     = true
  valid_redirect_uris = [
    "http://localhost:*"
  ]
  authorization {
    policy_enforcement_mode = "ENFORCING"
  }
}

resource "keycloak_role" "client_role_a" {
  realm_id  = keycloak_realm.this.id
  client_id = keycloak_openid_client.webapp.id
  name      = "demo.role.a"
}

resource "keycloak_role" "client_role_b" {
  realm_id  = keycloak_realm.this.id
  client_id = keycloak_openid_client.webapp.id
  name      = "demo.role.b"
}

resource "keycloak_openid_client" "webapp_public" {
  realm_id                     = keycloak_realm.this.id
  client_id                    = "${var.webapp_client_id}-public"
  client_secret                = "${var.webapp_client_id}-public"
  name                         = "Web Application Public"
  enabled                      = true
  access_type                  = "PUBLIC"
  standard_flow_enabled        = true
  implicit_flow_enabled        = false
  direct_access_grants_enabled = true
  service_accounts_enabled     = false
  valid_redirect_uris = [
    "http://localhost:*"
  ]
}

resource "keycloak_role" "client_role_c" {
  realm_id  = keycloak_realm.this.id
  client_id = keycloak_openid_client.webapp_public.id
  name      = "demo.role.c"
}

resource "keycloak_role" "client_role_d" {
  realm_id  = keycloak_realm.this.id
  client_id = keycloak_openid_client.webapp_public.id
  name      = "demo.role.d"
}
