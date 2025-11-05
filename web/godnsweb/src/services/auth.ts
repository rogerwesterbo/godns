// OIDC Configuration
export const OIDC_CONFIG = {
  authority: import.meta.env.VITE_KEYCLOAK_URL || 'http://localhost:14101',
  realm: import.meta.env.VITE_KEYCLOAK_REALM || 'godns',
  clientId: import.meta.env.VITE_KEYCLOAK_CLIENT_ID || 'godns-web',
  redirectUri: import.meta.env.VITE_REDIRECT_URI || 'http://localhost:14200/callback',
  postLogoutRedirectUri: import.meta.env.VITE_POST_LOGOUT_REDIRECT_URI || 'http://localhost:14200',
  scope: 'openid profile email',
};

// Storage keys
export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'godns_access_token',
  REFRESH_TOKEN: 'godns_refresh_token',
  ID_TOKEN: 'godns_id_token',
  EXPIRES_AT: 'godns_expires_at',
  CODE_VERIFIER: 'godns_code_verifier',
  STATE: 'godns_state',
};

export interface TokenResponse {
  access_token: string;
  refresh_token?: string;
  id_token?: string;
  expires_in: number;
  token_type: string;
}

export interface UserInfo {
  sub: string;
  email?: string;
  email_verified?: boolean;
  name?: string;
  preferred_username?: string;
  given_name?: string;
  family_name?: string;
  realm_access?: {
    roles: string[];
  };
}

/**
 * Generate a random string for PKCE code verifier
 */
export function generateCodeVerifier(): string {
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return base64UrlEncode(array);
}

/**
 * Generate code challenge from verifier using SHA-256
 */
export async function generateCodeChallenge(verifier: string): Promise<string> {
  const encoder = new TextEncoder();
  const data = encoder.encode(verifier);
  const hash = await crypto.subtle.digest('SHA-256', data);
  return base64UrlEncode(new Uint8Array(hash));
}

/**
 * Base64 URL encode without padding
 */
function base64UrlEncode(array: Uint8Array): string {
  const base64 = btoa(String.fromCharCode(...array));
  return base64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '');
}

/**
 * Generate random state parameter
 */
export function generateState(): string {
  const array = new Uint8Array(16);
  crypto.getRandomValues(array);
  return base64UrlEncode(array);
}

/**
 * Get OIDC endpoints
 */
export function getEndpoints() {
  const { authority, realm } = OIDC_CONFIG;
  const base = `${authority}/realms/${realm}/protocol/openid-connect`;

  return {
    authorization: `${base}/auth`,
    token: `${base}/token`,
    userinfo: `${base}/userinfo`,
    logout: `${base}/logout`,
    jwks: `${base}/certs`,
  };
}

/**
 * Build authorization URL with PKCE
 */
export async function buildAuthorizationUrl(): Promise<string> {
  const codeVerifier = generateCodeVerifier();
  const codeChallenge = await generateCodeChallenge(codeVerifier);
  const state = generateState();

  // Store for later use
  sessionStorage.setItem(STORAGE_KEYS.CODE_VERIFIER, codeVerifier);
  sessionStorage.setItem(STORAGE_KEYS.STATE, state);

  const params = new URLSearchParams({
    client_id: OIDC_CONFIG.clientId,
    redirect_uri: OIDC_CONFIG.redirectUri,
    response_type: 'code',
    scope: OIDC_CONFIG.scope,
    state: state,
    code_challenge: codeChallenge,
    code_challenge_method: 'S256',
  });

  return `${getEndpoints().authorization}?${params.toString()}`;
}

/**
 * Exchange authorization code for tokens
 */
export async function exchangeCodeForTokens(code: string, state: string): Promise<TokenResponse> {
  const savedState = sessionStorage.getItem(STORAGE_KEYS.STATE);
  const codeVerifier = sessionStorage.getItem(STORAGE_KEYS.CODE_VERIFIER);

  if (!savedState || savedState !== state) {
    throw new Error('Invalid state parameter');
  }

  if (!codeVerifier) {
    throw new Error('Code verifier not found');
  }

  const params = new URLSearchParams({
    grant_type: 'authorization_code',
    client_id: OIDC_CONFIG.clientId,
    redirect_uri: OIDC_CONFIG.redirectUri,
    code: code,
    code_verifier: codeVerifier,
  });

  const response = await fetch(getEndpoints().token, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: params.toString(),
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(`Token exchange failed: ${error}`);
  }

  const tokens: TokenResponse = await response.json();

  // Clean up session storage
  sessionStorage.removeItem(STORAGE_KEYS.STATE);
  sessionStorage.removeItem(STORAGE_KEYS.CODE_VERIFIER);

  return tokens;
}

/**
 * Refresh access token
 */
export async function refreshAccessToken(refreshToken: string): Promise<TokenResponse> {
  const params = new URLSearchParams({
    grant_type: 'refresh_token',
    client_id: OIDC_CONFIG.clientId,
    refresh_token: refreshToken,
  });

  const response = await fetch(getEndpoints().token, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body: params.toString(),
  });

  if (!response.ok) {
    throw new Error('Token refresh failed');
  }

  return await response.json();
}

/**
 * Get user info from userinfo endpoint
 */
export async function getUserInfo(accessToken: string): Promise<UserInfo> {
  const response = await fetch(getEndpoints().userinfo, {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });

  if (!response.ok) {
    throw new Error('Failed to fetch user info');
  }

  return await response.json();
}

/**
 * Logout and revoke tokens
 */
export async function logout(idToken?: string): Promise<void> {
  const params = new URLSearchParams({
    client_id: OIDC_CONFIG.clientId,
    post_logout_redirect_uri: OIDC_CONFIG.postLogoutRedirectUri,
  });

  if (idToken) {
    params.append('id_token_hint', idToken);
  }

  // Clear local storage
  Object.values(STORAGE_KEYS).forEach(key => {
    localStorage.removeItem(key);
    sessionStorage.removeItem(key);
  });

  // Redirect to Keycloak logout
  window.location.href = `${getEndpoints().logout}?${params.toString()}`;
}

/**
 * Store tokens in localStorage
 */
export function storeTokens(tokens: TokenResponse): void {
  const expiresAt = Date.now() + tokens.expires_in * 1000;

  localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, tokens.access_token);
  localStorage.setItem(STORAGE_KEYS.EXPIRES_AT, expiresAt.toString());

  if (tokens.refresh_token) {
    localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, tokens.refresh_token);
  }

  if (tokens.id_token) {
    localStorage.setItem(STORAGE_KEYS.ID_TOKEN, tokens.id_token);
  }
}

/**
 * Get stored access token
 */
export function getAccessToken(): string | null {
  return localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
}

/**
 * Get stored refresh token
 */
export function getRefreshToken(): string | null {
  return localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN);
}

/**
 * Get stored ID token
 */
export function getIdToken(): string | null {
  return localStorage.getItem(STORAGE_KEYS.ID_TOKEN);
}

/**
 * Check if access token is expired
 */
export function isTokenExpired(): boolean {
  const expiresAt = localStorage.getItem(STORAGE_KEYS.EXPIRES_AT);
  if (!expiresAt) return true;

  // Consider expired if less than 1 minute remaining
  return Date.now() >= parseInt(expiresAt) - 60000;
}

/**
 * Get valid access token, refreshing if needed
 */
export async function getValidAccessToken(): Promise<string | null> {
  const accessToken = getAccessToken();

  if (!accessToken) {
    return null;
  }

  if (!isTokenExpired()) {
    return accessToken;
  }

  // Try to refresh
  const refreshToken = getRefreshToken();
  if (!refreshToken) {
    return null;
  }

  try {
    const tokens = await refreshAccessToken(refreshToken);
    storeTokens(tokens);
    return tokens.access_token;
  } catch (error) {
    console.error('Token refresh failed:', error);
    return null;
  }
}
