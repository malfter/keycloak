resource "keycloak_user" "demo_user" {
  realm_id   = keycloak_realm.this.id
  username   = "demo"
  email      = "keycloak-demo+demo@alfter-web.de"
  first_name = "Testuser"
  last_name  = "Demo"
  enabled    = true
  initial_password {
    value     = "demo"
    temporary = false
  }
}

resource "keycloak_user_groups" "signer_groups" {
  realm_id = keycloak_realm.this.id
  user_id  = keycloak_user.demo_user.id

  group_ids = [
    keycloak_group.demo_group.id
  ]
}
