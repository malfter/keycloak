data "keycloak_openid_client_scope" "roles" {
  realm_id = keycloak_realm.this.id
  name     = "roles"
}

resource "keycloak_openid_user_client_role_protocol_mapper" "demo_roles" {
  realm_id        = keycloak_realm.this.id
  client_scope_id = data.keycloak_openid_client_scope.roles.id
  name            = "roles"
  claim_name      = "roles"
  multivalued     = true
  add_to_userinfo = true
}
