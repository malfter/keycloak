resource "keycloak_realm" "this" {
  realm   = var.realm_name
  enabled = true
}
