import { PUBLIC_API_URL } from '$env/static/public';
import ky from 'ky';

export const api = ky.create({
  prefixUrl: PUBLIC_API_URL || '/api/v1',
  hooks: {
    beforeRequest: [
      (request) => {
        if (typeof window !== 'undefined') {
          const token = localStorage.getItem('token');
          if (token) {
            request.headers.set('Authorization', `Bearer ${token}`);
          }
        }
      }
    ],
    afterResponse: [
      (_request, _options, response) => {
        if (response.status === 401) {
          if (typeof window !== 'undefined') {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/login';
          }
        }
      }
    ]
  }
});

// Types
export interface User {
  id: string;
  email: string;
  name: string;
  email_verified?: boolean;
  created_at: string;
}

export interface Organization {
  id: string;
  name: string;
  owner_id: string;
  created_at: string;
}

export interface Application {
  id: string;
  organization_id: string;
  name: string;
  description: string;
  created_at: string;
}

export interface ApplicationSecret {
  id: string;
  application_id: string;
  key: string;
  name: string;
  created_at: string;
}

export interface EventType {
  id: string;
  application_id: string;
  name: string;
  description: string;
  schema?: object;
  created_at: string;
}

export interface Subscription {
  id: string;
  application_id: string;
  event_types: string[];
  target_url: string;
  description: string;
  enabled: boolean;
  created_at: string;
}

export interface Event {
  id: string;
  application_id: string;
  event_type: string;
  payload: object;
  metadata?: object;
  created_at: string;
}

export interface DeliveryAttempt {
  id: string;
  event_id: string;
  subscription_id: string;
  status: string;
  status_code: number;
  request_body: string;
  response_body: string;
  duration_ms: number;
  attempt_number: number;
  created_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface StatsResponse {
  total_events: number;
  total_deliveries: number;
  success_rate: number;
  total_subscriptions: number;
  total_applications: number;
}

export interface ChartDataPoint {
  date: string;
  success: number;
  failed: number;
  total: number;
}

// API functions
export const authApi = {
  register: (data: { email: string; password: string; name: string }) =>
    api.post('auth/register', { json: data }).json<{ token: string; refresh_token: string; user: User }>(),
  login: (data: { email: string; password: string }) =>
    api.post('auth/login', { json: data }).json<{ token: string; refresh_token: string; user: User }>(),
  refresh: (refresh_token: string) =>
    api.post('auth/refresh', { json: { refresh_token } }).json<{ token: string }>(),
  sendVerification: () =>
    api.post('auth/send-verification').json<{ message: string; token?: string }>(),
  verifyEmail: (token: string) =>
    api.post(`auth/verify/${token}`).json<{ message: string }>(),
};

export const orgApi = {
  list: () => api.get('organizations').json<Organization[]>(),
  get: (id: string) => api.get(`organizations/${id}`).json<Organization>(),
  create: (data: { name: string }) =>
    api.post('organizations', { json: data }).json<Organization>(),
  update: (id: string, data: { name: string }) =>
    api.put(`organizations/${id}`, { json: data }).json<Organization>(),
  delete: (id: string) => api.delete(`organizations/${id}`).json(),
  dashboard: (id: string) => api.get(`organizations/${id}/dashboard`).json<StatsResponse>(),
  members: (id: string) => api.get(`organizations/${id}/members`).json<any[]>(),
  inviteMember: (id: string, data: { email: string; role: string }) =>
    api.post(`organizations/${id}/members`, { json: data }).json(),
  removeMember: (orgId: string, userId: string) =>
    api.delete(`organizations/${orgId}/members/${userId}`).json(),
  acceptInvitation: (token: string) =>
    api.post(`invitations/${token}/accept`).json<any>(),
};

export const appApi = {
  list: (orgId: string) => api.get(`organizations/${orgId}/applications`).json<Application[]>(),
  get: (id: string) => api.get(`applications/${id}`).json<Application>(),
  create: (orgId: string, data: { name: string; description: string }) =>
    api.post(`organizations/${orgId}/applications`, { json: data }).json<Application>(),
  update: (id: string, data: { name: string; description: string }) =>
    api.put(`applications/${id}`, { json: data }).json<Application>(),
  delete: (id: string) => api.delete(`applications/${id}`).json(),
  dashboard: (id: string) => api.get(`applications/${id}/dashboard`).json<StatsResponse & { chart: ChartDataPoint[] }>(),
};

export const secretApi = {
  list: (appId: string, params?: { page?: number; per_page?: number }) =>
    api.get(`applications/${appId}/secrets`, { searchParams: params as any }).json<PaginatedResponse<ApplicationSecret>>(),
  create: (appId: string, data: { name: string }) =>
    api.post(`applications/${appId}/secrets`, { json: data }).json<ApplicationSecret>(),
  delete: (appId: string, secretId: string) =>
    api.delete(`applications/${appId}/secrets/${secretId}`).json(),
};

export const eventTypeApi = {
  list: (appId: string, params?: { page?: number; per_page?: number }) =>
    api.get(`applications/${appId}/event-types`, { searchParams: params as any }).json<PaginatedResponse<EventType>>(),
  create: (appId: string, data: { name: string; description: string; schema?: object }) =>
    api.post(`applications/${appId}/event-types`, { json: data }).json<EventType>(),
  delete: (appId: string, etId: string) =>
    api.delete(`applications/${appId}/event-types/${etId}`).json(),
};

export const subscriptionApi = {
  list: (appId: string, params?: { page?: number; per_page?: number }) =>
    api.get(`applications/${appId}/subscriptions`, { searchParams: params as any }).json<PaginatedResponse<Subscription>>(),
  get: (id: string) => api.get(`subscriptions/${id}`).json<Subscription>(),
  create: (appId: string, data: { event_types: string[]; target_url: string; description: string }) =>
    api.post(`applications/${appId}/subscriptions`, { json: data }).json<Subscription>(),
  update: (id: string, data: Partial<Subscription>) =>
    api.put(`subscriptions/${id}`, { json: data }).json<Subscription>(),
  delete: (id: string) => api.delete(`subscriptions/${id}`).json(),
  test: (id: string) => api.post(`subscriptions/${id}/test`).json(),
};

export const eventApi = {
  send: (appId: string, data: { event_type: string; payload: object; metadata?: object }) =>
    api.post(`applications/${appId}/events`, { json: data }).json<Event>(),
  list: (appId: string, params?: { page?: number; per_page?: number }) =>
    api.get(`applications/${appId}/events`, { searchParams: params as any }).json<PaginatedResponse<Event>>(),
  get: (id: string) => api.get(`events/${id}`).json<Event>(),
};

export const deliveryApi = {
  list: (appId: string, params?: { page?: number; per_page?: number; status?: string }) =>
    api.get(`applications/${appId}/deliveries`, { searchParams: params as any }).json<PaginatedResponse<DeliveryAttempt>>(),
  get: (id: string) => api.get(`deliveries/${id}`).json<DeliveryAttempt>(),
  retry: (id: string) => api.post(`deliveries/${id}/retry`).json(),
};
