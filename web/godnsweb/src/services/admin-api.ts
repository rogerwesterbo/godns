import { getValidAccessToken } from './auth';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:14000';

export interface SystemStats {
  cache: CacheStats;
  rate_limiter: RateLimiterStats;
  load_balancer: LoadBalancerStats;
  health_check: HealthCheckStats;
  query_log: QueryLogStats;
}

export interface CacheStats {
  enabled: boolean;
  current_size: number;
  max_size: number;
  ttl_minutes: number;
  hit_rate?: string;
}

export interface CacheStatsDetailed {
  enabled: boolean;
  size: number;
  capacity: number;
  hits: number;
  misses: number;
  hit_rate: number | string; // Decimal 0-1 (e.g., 0.931 = 93.1%) or string "93.10%"
  evictions: number;
}

export interface BackendHealth {
  name: string;
  type: string;
  value: string;
  address: string;
  weight: number;
  healthy: boolean;
  enabled: boolean;
  response_time_ms: number;
  last_check: string;
}

export interface LoadBalancerStats {
  enabled: boolean;
  strategy: string;
  backend_groups: number;
  total_backends: number;
  healthy_backends: number;
  backends?: BackendHealth[];
}

export interface HealthCheckTarget {
  target: string;
  host: string;
  port: number;
  protocol: string;
  healthy: boolean;
  last_check: string;
  response_time_ms: number;
  last_error: string;
}

export interface HealthCheckStats {
  enabled: boolean;
  total_targets: number;
  healthy_targets: number;
  interval_seconds: number;
  timeout_seconds: number;
  targets?: HealthCheckTarget[];
  results?: HealthCheckTarget[];
}

export interface QueryLogStats {
  enabled: boolean;
  total_queries: number;
  cached_queries: number;
  blocked_queries: number;
  cache_hit_rate: number; // Decimal 0-1 (e.g., 0.931 = 93.1%)
}

export interface RateLimiterStats {
  enabled: boolean;
  qps: number;
  burst: number;
  active_limiters: number;
  total_blocked: number;
}

class AdminApiError extends Error {
  status: number;
  response?: unknown;

  constructor(status: number, message: string, response?: unknown) {
    super(message);
    this.name = 'AdminApiError';
    this.status = status;
    this.response = response;
  }
}

async function adminApiRequest<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = await getValidAccessToken();

  if (!token) {
    throw new AdminApiError(401, 'Not authenticated');
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
    let errorMessage = `Admin API request failed: ${response.statusText}`;
    let errorData: unknown;

    try {
      errorData = await response.json();
      if (errorData && typeof errorData === 'object' && 'error' in errorData) {
        errorMessage = String(errorData.error);
      }
    } catch {
      // Ignore JSON parse errors
    }

    throw new AdminApiError(response.status, errorMessage, errorData);
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

// Admin endpoints
export async function getSystemStats(): Promise<SystemStats> {
  return adminApiRequest<SystemStats>('/api/v1/admin/stats');
}

export async function getCacheStats(): Promise<CacheStatsDetailed> {
  return adminApiRequest<CacheStatsDetailed>('/api/v1/admin/cache/stats');
}

export async function clearCache(): Promise<void> {
  return adminApiRequest<void>('/api/v1/admin/cache/clear', {
    method: 'POST',
  });
}

export async function getLoadBalancerStats(): Promise<LoadBalancerStats> {
  return adminApiRequest<LoadBalancerStats>('/api/v1/admin/loadbalancer/stats');
}

export async function getHealthCheckStats(): Promise<HealthCheckStats> {
  return adminApiRequest<HealthCheckStats>('/api/v1/admin/healthcheck/stats');
}

export async function getQueryLogStats(): Promise<QueryLogStats> {
  return adminApiRequest<QueryLogStats>('/api/v1/admin/querylog/stats');
}

export async function getRateLimiterStats(): Promise<RateLimiterStats> {
  return adminApiRequest<RateLimiterStats>('/api/v1/admin/ratelimiter/stats');
}

export { AdminApiError };
