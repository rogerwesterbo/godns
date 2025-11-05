# GoDNS Web UI

A modern React-based web interface for the GoDNS DNS management system. Built with Vite, TypeScript, and Radix UI.

[![React](https://img.shields.io/badge/React-19-blue)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9-blue)](https://www.typescriptlang.org/)
[![Vite](https://img.shields.io/badge/Vite-Rolldown-purple)](https://vite.dev/)

## ✨ Features

- 🔐 **OAuth 2.0 + PKCE Authentication** - Secure Keycloak integration
- 🌓 **Dark/Light Theme** - Persistent theme switching with Radix UI
- 🔍 **Real-time Search** - Instant zone and record search
- 📊 **Interactive Dashboard** - Live metrics and statistics
- 🌐 **Full CRUD Operations** - Create, read, update, delete zones and records
- 🔄 **Auto Token Refresh** - Seamless session management
- 📱 **Responsive Design** - Mobile-optimized interface
- ⚡ **Fast Performance** - Vite with Rolldown bundler

## Quick Start

### Prerequisites

- Node.js 20.x or later
- npm or yarn
- GoDNS API server running (http://localhost:14000)
- Keycloak server running (http://localhost:14101)

### Installation

```bash
# Install dependencies
npm install

# Start development server (runs on port 14200)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

The application will be available at http://localhost:14200

### Docker Deployment

```bash
# Build Docker image
docker build -t godnsweb:latest .

# Run container
docker run -p 8080:8080 \
  -e VITE_KEYCLOAK_URL=http://keycloak:8080 \
  -e VITE_API_BASE_URL=http://godns-api:8080 \
  godnsweb:latest
```

See [docs/DOCKER.md](./docs/DOCKER.md) for detailed Docker and Kubernetes deployment.

### Kubernetes/Helm

```bash
# Install with Helm
helm install godnsweb oci://ghcr.io/rogerwesterbo/helm/godnsweb \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=godns.example.com
```

See [charts/godnsweb/README.md](../../charts/godnsweb/README.md) for Helm chart documentation.

## Tech Stack

- **[React 19](https://react.dev/)** - Latest React with new features
- **[TypeScript 5.9](https://www.typescriptlang.org/)** - Type-safe development
- **[Vite](https://vite.dev/)** (Rolldown 7.1.20) - Lightning-fast build tool
- **[Radix UI Themes](https://www.radix-ui.com/themes)** - Accessible components
- **[React Router 7](https://reactrouter.com/)** - Client-side routing
- **[jwt-decode](https://github.com/auth0/jwt-decode)** - JWT token handling

## Documentation

- **[Features Guide](./docs/FEATURES.md)** - Detailed feature documentation
- **[Authentication Setup](./docs/AUTH_SETUP.md)** - OAuth 2.0 PKCE configuration
- **[Components Guide](./docs/COMPONENTS.md)** - Component architecture
- **[Routes & Navigation](./docs/ROUTES.md)** - Routing configuration
- **[Theme System](./docs/THEME_GUIDE.md)** - Theme customization
- **[Docker Deployment](./docs/DOCKER.md)** - Container and Kubernetes setup
- **[Quick Start](./docs/QUICK_START.md)** - Get started in 5 minutes

## Project Structure

```
web/godnsweb/
├── src/
│   ├── components/      # Reusable UI components
│   │   ├── Header.tsx
│   │   ├── Sidebar.tsx
│   │   ├── SearchBar.tsx
│   │   ├── CreateZoneDialog.tsx
│   │   ├── RecordDialog.tsx
│   │   └── ProtectedRoute.tsx
│   ├── contexts/        # React contexts
│   │   ├── AuthContext.ts
│   │   ├── AuthProvider.tsx
│   │   ├── ThemeContext.ts
│   │   └── ThemeProvider.tsx
│   ├── pages/           # Application pages
│   │   ├── DashboardPage.tsx
│   │   ├── ZonesPage.tsx
│   │   ├── ZoneDetailPage.tsx
│   │   ├── RecordsPage.tsx
│   │   ├── SearchPage.tsx
│   │   ├── ProfilePage.tsx
│   │   ├── LoginPage.tsx
│   │   ├── CallbackPage.tsx
│   │   └── error pages...
│   ├── services/        # API and auth services
│   │   ├── api.ts       # API client with all endpoints
│   │   └── auth.ts      # OAuth 2.0 PKCE implementation
│   ├── App.tsx          # App router & providers
│   └── main.tsx         # Application entry point
├── docs/                # Documentation
├── Dockerfile           # Multi-stage Docker build
├── .dockerignore        # Docker build exclusions
└── package.json         # Dependencies and scripts
```

## Environment Variables

Create a `.env` file in the root directory:

```bash
# Keycloak/OIDC Configuration
VITE_KEYCLOAK_URL=http://localhost:14101
VITE_KEYCLOAK_REALM=godns
VITE_KEYCLOAK_CLIENT_ID=godns-web
VITE_REDIRECT_URI=http://localhost:14200/callback
VITE_POST_LOGOUT_REDIRECT_URI=http://localhost:14200

# API Configuration
VITE_API_BASE_URL=http://localhost:14000
```

## Development

```bash
# Start dev server with HMR
npm run dev

# Type checking
npm run build

# Linting
npm run lint

# Format code (if prettier is configured)
npm run format
```

## Building for Production

```bash
# Build optimized bundle
npm run build

# Preview production build locally
npm run preview

# Build Docker image
docker build -t godnsweb:latest .
```

## Security

### Authentication

- **OAuth 2.0 with PKCE** - Industry-standard authorization
- **JWT Token Management** - Secure token storage and refresh
- **Auto Token Refresh** - Seamless session renewal
- **Secure Logout** - Proper token revocation

### Container Security

- **Google Distroless Base** - Minimal attack surface
- **Non-root User** - Runs as UID 65532
- **Read-only Filesystem** - Enhanced security
- **No Shell/Package Manager** - Reduced vulnerability

### Code Security

- **TypeScript** - Type safety
- **Security Headers** - X-Frame-Options, CSP, etc.
- **Input Validation** - Client-side validation
- **CORS Protection** - Configured origins

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

See the [LICENSE](../../LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/rogerwesterbo/godns/issues)
- **Documentation**: [Main Docs](../../docs/)
- **API Docs**: [API Documentation](../../docs/API_DOCUMENTATION.md)
