resource "keycloak_group" "demo_group" {
  realm_id = keycloak_realm.this.id
  name     = "demo-group"
}

resource "keycloak_group_roles" "group_roles" {
  realm_id = keycloak_realm.this.id
  group_id = keycloak_group.demo_group.id

  role_ids = [
    keycloak_role.client_role_a.id,
    keycloak_role.client_role_b.id,
    keycloak_role.client_role_c.id,
    keycloak_role.client_role_d.id,
  ]
}
