#!/usr/bin/env bash

: "${OAUTH2_USERNAME:="demo"}"
: "${OAUTH2_PASSWORD:="demo"}"

: "${OAUTH2_CLIENT_ID:="demo-webapp-public"}"
: "${OAUTH2_URL_TOKEN:="http://localhost:9080/realms/keycloak-demo/protocol/openid-connect/token"}"
: "${OAUTH2_URL_USERINFO:="http://localhost:9080/realms/keycloak-demo/protocol/openid-connect/userinfo"}"

# Realm informations
# http://localhost:9080/realms/keycloak-demo/.well-known/openid-configuration

RESPONSE_TOKEN=$(curl -s -X POST "${OAUTH2_URL_TOKEN}" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=${OAUTH2_CLIENT_ID}" \
  -d "grant_type=password" \
  -d "username=${OAUTH2_USERNAME}" \
  -d "password=${OAUTH2_PASSWORD}" \
  -d "scope=openid profile email offline_access")

echo "RESPONSE TOKEN:"
echo "${RESPONSE_TOKEN}" | jq .
echo ""

ACCESS_TOKEN=$(echo "${RESPONSE_TOKEN}" | jq -r .access_token)
ID_TOKEN=$(echo "${RESPONSE_TOKEN}" | jq -r .id_token)

RESPONSE_USERINFO=$(curl -s -X GET "${OAUTH2_URL_USERINFO}" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}")

echo "RESPONSE USERINFO:"
echo "${RESPONSE_USERINFO}" | jq .
echo ""
