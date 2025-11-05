import { getValidAccessToken } from './auth';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:14000';

export interface DNSRecord {
  name: string;
  type: string;
  ttl: number;
  value: string;
}

export interface DNSZone {
  domain: string;
  records: DNSRecord[];
}

export interface SearchResult {
  type: 'zone' | 'record';
  zone?: DNSZone;
  record?: DNSRecord & { zone: string };
  highlight?: string;
}

export interface SearchResponse {
  results: SearchResult[];
  total: number;
  query: string;
}

class ApiError extends Error {
  status: number;
  response?: unknown;

  constructor(status: number, message: string, response?: unknown) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.response = response;
  }
}

async function apiRequest<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = await getValidAccessToken();

  if (!token) {
    throw new ApiError(401, 'Not authenticated');
  }

  const url = `${API_BASE_URL}${endpoint}`;

  const response = await fetch(url, {
    ...options,
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    let errorMessage = `API request failed: ${response.statusText}`;
    let errorData: unknown;

    try {
      errorData = await response.json();
      if (errorData && typeof errorData === 'object' && 'error' in errorData) {
        errorMessage = String(errorData.error);
      }
    } catch {
      // Ignore JSON parse errors
    }

    throw new ApiError(response.status, errorMessage, errorData);
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

// Zone endpoints
export async function listZones(): Promise<DNSZone[]> {
  return apiRequest<DNSZone[]>('/api/v1/zones');
}

export async function getZone(domain: string): Promise<DNSZone> {
  return apiRequest<DNSZone>(`/api/v1/zones/${encodeURIComponent(domain)}`);
}

export async function createZone(zone: DNSZone): Promise<DNSZone> {
  return apiRequest<DNSZone>('/api/v1/zones', {
    method: 'POST',
    body: JSON.stringify(zone),
  });
}

export async function updateZone(domain: string, zone: DNSZone): Promise<DNSZone> {
  return apiRequest<DNSZone>(`/api/v1/zones/${encodeURIComponent(domain)}`, {
    method: 'PUT',
    body: JSON.stringify(zone),
  });
}

export async function deleteZone(domain: string): Promise<void> {
  return apiRequest<void>(`/api/v1/zones/${encodeURIComponent(domain)}`, {
    method: 'DELETE',
  });
}

// Record endpoints
export async function createRecord(domain: string, record: DNSRecord): Promise<DNSRecord> {
  return apiRequest<DNSRecord>(`/api/v1/zones/${encodeURIComponent(domain)}/records`, {
    method: 'POST',
    body: JSON.stringify(record),
  });
}

export async function getRecord(domain: string, name: string, type: string): Promise<DNSRecord> {
  return apiRequest<DNSRecord>(
    `/api/v1/zones/${encodeURIComponent(domain)}/records/${encodeURIComponent(name)}/${encodeURIComponent(type)}`
  );
}

export async function updateRecord(
  domain: string,
  name: string,
  type: string,
  record: DNSRecord
): Promise<DNSRecord> {
  return apiRequest<DNSRecord>(
    `/api/v1/zones/${encodeURIComponent(domain)}/records/${encodeURIComponent(name)}/${encodeURIComponent(type)}`,
    {
      method: 'PUT',
      body: JSON.stringify(record),
    }
  );
}

export async function deleteRecord(domain: string, name: string, type: string): Promise<void> {
  return apiRequest<void>(
    `/api/v1/zones/${encodeURIComponent(domain)}/records/${encodeURIComponent(name)}/${encodeURIComponent(type)}`,
    {
      method: 'DELETE',
    }
  );
}

// Search endpoint
export async function search(
  query: string,
  types?: ('zone' | 'record')[]
): Promise<SearchResponse> {
  const params = new URLSearchParams({ q: query });

  if (types && types.length > 0) {
    types.forEach(type => params.append('type', type));
  }

  return apiRequest<SearchResponse>(`/api/v1/search?${params.toString()}`);
}

export { ApiError };
