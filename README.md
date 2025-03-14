# 🔐 keycloak

<img src="assets/Keycloak_Logo.png" alt="Keycloak Logo" width="100">

This project is used to test some [Keycloak](https://www.keycloak.org/) configurations 🛠️.

## 📖 Table of Contents

- [🔐 keycloak](#-keycloak)
  - [📖 Table of Contents](#-table-of-contents)
  - [📌 Requirements](#-requirements)
  - [🆔 OIDC Example](#-oidc-example)
    - [🌱 Get Started](#-get-started)
  - [🔗 Further Links](#-further-links)

## 📌 Requirements

> ℹ️ If you don't want to install anything locally, you can also use the [devcontainer](.devcontainer/devcontainer.json) environment, which only requires a container runtime such as [podman](https://podman.io/)/[docker](https://docker.com).

To work with this project, you need to install some dependencies:

- [https://go.dev/] (Demo apps are written in Go)
- [https://opentofu.org/] (OpenTofu is used to configure Keycloak)

## 🆔 OIDC Example

This example shows the use of Keycloak as an identity provider in combination with a web app and a backend service.

1. Opening the web app
   - [http://localhost:9081/]
2. In order to use the web app, a login is required ("click Login").
   - [http://localhost:9081/login]
3. After logging in, you will be returned to the web app (redirect to `/callback`).
   - [http://localhost:9081/callback]
4. The information retrieved from the backend service is now displayed.

### 🌱 Get Started

```bash
# Start keycloak instance
docker compose up

# Configure keycloak with terraform
cd config/keycloak
tofu init
tofu apply

# Start backend
cd demo-apps/demo-backend
go run .

# Start webapp
cd demo-apps/demo-webapp
go run .

# Open webapp and "click" login
# Username: demo
# Password: demo
open http://localhost:9081/
```

## 🔗 Further Links

- Keycloak Documentation
  - [https://www.keycloak.org/documentation]
  - [https://www.keycloak.org/server/containers]
- Keycloak OpenTofu Provider
  - [https://search.opentofu.org/provider/keycloak/keycloak/]
