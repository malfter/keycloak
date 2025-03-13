output "webapp_client_secret" {
  value     = keycloak_openid_client.webapp.client_secret
  sensitive = true
}
