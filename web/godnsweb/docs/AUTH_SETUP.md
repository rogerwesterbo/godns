# OAuth 2.0 / OIDC Authentication Setup

This guide explains how to set up OpenID Connect authentication for the GoDNS web application using Keycloak with PKCE.

## 🔐 Authentication Flow

The web application uses the **Authorization Code Flow with PKCE (Proof Key for Code Exchange)** which is the recommended OAuth 2.0 flow for SPAs (Single Page Applications).

### Flow Diagram

```
1. User clicks "Login" → Redirected to Keycloak
2. User authenticates with Keycloak
3. Keycloak redirects back with authorization code
4. App exchanges code for tokens (using PKCE)
5. Tokens stored in localStorage
6. User is authenticated and can access protected routes
```

## ⚙️ Keycloak Configuration

### Step 1: Create the Web Client

You need to create a new client in Keycloak for the web application.

1. Open Keycloak Admin Console: http://localhost:14101
2. Login with admin credentials (admin/admin)
3. Select the `godns` realm
4. Navigate to **Clients** → **Create client**

### Step 2: Client Configuration

**General Settings:**
- Client type: `OpenID Connect`
- Client ID: `godns-web`
- Name: `GoDNS Web Application`
- Click **Next**

**Capability config:**
- ✅ Client authentication: **OFF** (Public client)
- ✅ Authorization: **OFF**
- ✅ Authentication flow:
  - ✅ Standard flow: **ON**
  - ❌ Direct access grants: **OFF**
  - ❌ Implicit flow: **OFF**
  - ❌ Service accounts roles: **OFF**
- Click **Next**

**Login settings:**
- Root URL: `http://localhost:14200`
- Home URL: `http://localhost:14200`
- Valid redirect URIs: `http://localhost:14200/callback`
- Valid post logout redirect URIs: `http://localhost:14200`
- Web origins: `http://localhost:14200`
- Click **Save**

### Step 3: Advanced Settings (Optional)

Go to the **Advanced** tab and configure:

- **Proof Key for Code Exchange Code Challenge Method**: `S256` (should be enabled by default for public clients)
- **Access Token Lifespan**: `3600` seconds (1 hour)

## 🔧 Environment Configuration

Create a `.env` file in `web/godnsweb/`:

```bash
VITE_KEYCLOAK_URL=http://localhost:14101
VITE_KEYCLOAK_REALM=godns
VITE_KEYCLOAK_CLIENT_ID=godns-web
VITE_REDIRECT_URI=http://localhost:14200/callback
VITE_POST_LOGOUT_REDIRECT_URI=http://localhost:14200
```

## 🚀 Running the Application

```bash
cd web/godnsweb
npm install
npm run dev
```

Visit: http://localhost:14200

## 📝 How It Works

### 1. Login Flow

When a user clicks "Login":

```typescript
// 1. Generate PKCE code verifier and challenge
const codeVerifier = generateCodeVerifier();  // Random string
const codeChallenge = await generateCodeChallenge(codeVerifier);  // SHA-256 hash

// 2. Store verifier in session
sessionStorage.setItem('code_verifier', codeVerifier);

// 3. Redirect to Keycloak with challenge
window.location.href = `${keycloakURL}/auth?
  client_id=godns-web&
  redirect_uri=http://localhost:14200/callback&
  response_type=code&
  scope=openid profile email&
  code_challenge=${codeChallenge}&
  code_challenge_method=S256&
  state=${randomState}`;
```

### 2. Callback Handling

After Keycloak authentication, user is redirected to `/callback`:

```typescript
// 1. Extract code and state from URL
const code = searchParams.get('code');
const state = searchParams.get('state');

// 2. Verify state matches
if (savedState !== state) throw new Error('Invalid state');

// 3. Exchange code for tokens with PKCE verifier
const tokens = await fetch(tokenEndpoint, {
  method: 'POST',
  body: new URLSearchParams({
    grant_type: 'authorization_code',
    client_id: 'godns-web',
    redirect_uri: 'http://localhost:14200/callback',
    code: code,
    code_verifier: savedCodeVerifier  // Proves we initiated the request
  })
});

// 4. Store tokens
localStorage.setItem('access_token', tokens.access_token);
localStorage.setItem('refresh_token', tokens.refresh_token);

// 5. Redirect to dashboard
navigate('/');
```

### 3. Protected Routes

Routes are protected using the `ProtectedRoute` component:

```typescript
<Route path="/" element={
  <ProtectedRoute>
    <Layout><DashboardPage /></Layout>
  </ProtectedRoute>
} />
```

If not authenticated, user is redirected to `/login`.

### 4. Token Refresh

Tokens are automatically refreshed when expired:

```typescript
export async function getValidAccessToken() {
  const accessToken = getAccessToken();
  
  if (!isTokenExpired()) {
    return accessToken;
  }
  
  // Auto-refresh
  const refreshToken = getRefreshToken();
  const newTokens = await refreshAccessToken(refreshToken);
  storeTokens(newTokens);
  
  return newTokens.access_token;
}
```

### 5. Logout

Logout clears tokens and redirects to Keycloak logout:

```typescript
export async function logout() {
  // Clear local storage
  localStorage.removeItem('access_token');
  localStorage.removeItem('refresh_token');
  
  // Redirect to Keycloak logout
  window.location.href = `${keycloakURL}/logout?
    client_id=godns-web&
    post_logout_redirect_uri=http://localhost:14200&
    id_token_hint=${idToken}`;
}
```

## 🔑 User Information

User information is extracted from the JWT token:

```typescript
import { jwtDecode } from 'jwt-decode';

const token = getAccessToken();
const user = jwtDecode(token);

// Available fields:
// - sub: User ID
// - email: Email address
// - preferred_username: Username
// - name: Full name
// - given_name: First name
// - family_name: Last name
// - realm_access.roles: Array of roles
```

The ProfilePage component displays this information.

## 🛡️ Security Features

### PKCE (Proof Key for Code Exchange)

PKCE prevents authorization code interception attacks:

1. **Code Verifier**: Random string (43-128 characters)
2. **Code Challenge**: SHA-256 hash of verifier
3. Challenge sent to Keycloak during authorization
4. Verifier sent during token exchange
5. Keycloak verifies hash(verifier) == challenge

This ensures only the app that initiated the flow can exchange the code.

### State Parameter

Random state parameter prevents CSRF attacks:

```typescript
const state = generateRandomString();
sessionStorage.setItem('state', state);

// After callback, verify:
if (receivedState !== savedState) {
  throw new Error('CSRF attempt detected');
}
```

### Token Storage

- **Access Token**: Short-lived (1 hour), stored in localStorage
- **Refresh Token**: Longer-lived, used to get new access tokens
- **ID Token**: Contains user information

### Secure Headers

When calling the API, include the access token:

```typescript
const token = await getValidAccessToken();

fetch('/api/v1/zones', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

## 📦 API Integration

Create an API client that automatically includes auth:

```typescript
// src/services/api.ts
import { getValidAccessToken } from './auth';

export async function apiClient(url: string, options = {}) {
  const token = await getValidAccessToken();
  
  if (!token) {
    throw new Error('Not authenticated');
  }
  
  return fetch(`http://localhost:14000${url}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });
}

// Usage:
const response = await apiClient('/api/v1/zones');
const zones = await response.json();
```

## 🧪 Testing

### Test User

The default test user (created by init script):

- **Username**: `testuser`
- **Password**: `password`
- **Email**: `testuser@godns.local`
- **Role**: `dns-admin`

### Manual Testing

1. Start services: `docker-compose up -d`
2. Start web app: `cd web/godnsweb && npm run dev`
3. Open: http://localhost:14200
4. Click "Login"
5. Login with testuser/password
6. You should be redirected to dashboard
7. Check Profile page for user info

### Testing Logout

1. Click user avatar in header
2. Click "Logout"
3. You should be logged out and redirected to login page

### Testing Token Refresh

Tokens auto-refresh 1 minute before expiry. To test:

1. Login
2. Wait for token to expire (or manually adjust expiry in localStorage)
3. Make an API call
4. Token should auto-refresh

## 🐛 Troubleshooting

### "Invalid redirect URI"

- Ensure `http://localhost:14200/callback` is in Keycloak's "Valid redirect URIs"
- Check for trailing slashes
- Verify port numbers match

### "CORS error"

Update GoDNS API CORS settings in `.env`:

```bash
HTTP_API_CORS_ALLOWED_ORIGINS=http://localhost:14000,http://localhost:14200
```

### "Token verification failed"

- Verify Keycloak URL in `.env` matches running Keycloak
- Check realm name is correct (`godns`)
- Ensure client ID is `godns-web`

### "Code challenge failed"

- Ensure PKCE is enabled in Keycloak client settings
- Verify `code_challenge_method` is `S256`
- Check browser supports crypto.subtle API

### "State mismatch"

- Clear browser cache and cookies
- Check sessionStorage is working
- Ensure state parameter is being generated and stored

## 📚 References

- [OAuth 2.0 PKCE RFC](https://datatracker.ietf.org/doc/html/rfc7636)
- [OpenID Connect Core](https://openid.net/specs/openid-connect-core-1_0.html)
- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [OWASP OAuth 2.0 Security](https://cheatsheetseries.owasp.org/cheatsheets/OAuth2_Cheat_Sheet.html)

## 🔄 Production Checklist

- [ ] Use HTTPS for all URLs
- [ ] Configure proper CORS origins
- [ ] Set secure session storage
- [ ] Implement token rotation
- [ ] Add rate limiting
- [ ] Monitor authentication logs
- [ ] Implement MFA
- [ ] Use secure cookie flags
- [ ] Add CSP headers
- [ ] Regular security audits
