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

// @contact.name GoDNS Support
// @contact.url https://github.com/rogerwesterbo/godns
// @contact.email support@example.com

// @license.name MIT
// @license.url https://github.com/rogerwesterbo/godns/blob/main/LICENSE

// @host localhost:14082
// @BasePath /

// @schemes http https
// @produce json
// @consumes json

// @tag.name Health
// @tag.description Health check and readiness endpoints

// @tag.name Zones
// @tag.description DNS zone management operations

// @tag.name Records
// @tag.description DNS record management operations
