package main

// @title GoDNS API
// @version 1.0
// @description RESTful API for managing DNS zones and records in GoDNS
// @description
// @description This API provides CRUD operations for DNS zones and records, storing data in Valkey (Redis-compatible).
// @description
// @description ## Features
// @description - Create, read, update, and delete DNS zones
// @description - Manage individual DNS records within zones
// @description - Support for various record types (A, AAAA, CNAME, MX, NS, TXT, PTR, SRV, SOA, CAA)
// @description - Health and readiness endpoints for Kubernetes deployments
// @description
// @description ## Authentication
// @description This API uses OAuth2/JWT authentication via Keycloak.
// @description
// @description To authenticate:
// @description 1. Click "Authorize" button
// @description 2. For BearerAuth: Enter "Bearer YOUR_JWT_TOKEN"
// @description 3. For OAuth2: Enter credentials (testuser/password for development)
// @description
// @description To get a token manually (HTTP):
// @description ```
// @description curl -X POST "http://localhost:14101/realms/godns/protocol/openid-connect/token" \
// @description   -H "Content-Type: application/x-www-form-urlencoded" \
// @description   -d "client_id=godns-cli" \
// @description   -d "username=testuser" \
// @description   -d "password=password" \
// @description   -d "grant_type=password"
// @description ```
// @description Or use HTTPS on port 14102 for production.

// @contact.name GoDNS Support
// @contact.url https://github.com/rogerwesterbo/godns
// @contact.email support@example.com

// @license.name MIT
// @license.url https://github.com/rogerwesterbo/godns/blob/main/LICENSE

// @host localhost:14000
// @BasePath /

// @schemes http https
// @produce json
// @consumes json

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl http://localhost:14101/realms/godns/protocol/openid-connect/token

// @tag.name Health
// @tag.description Health check and readiness endpoints

// @tag.name Zones
// @tag.description DNS zone management operations

// @tag.name Records
// @tag.description DNS record management operations
