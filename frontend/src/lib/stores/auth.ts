import { writable } from 'svelte/store';
import type { User } from '$lib/api/client';

function createAuthStore() {
  const stored = typeof window !== 'undefined' ? localStorage.getItem('user') : null;
  const initial = stored ? JSON.parse(stored) : null;

  const { subscribe, set, update } = writable<{
    user: User | null;
    token: string | null;
    loading: boolean;
  }>({
    user: initial,
    token: typeof window !== 'undefined' ? localStorage.getItem('token') : null,
    loading: false
  });

  return {
    subscribe,
    login: (user: User, token: string) => {
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      set({ user, token, loading: false });
    },
    logout: () => {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      set({ user: null, token: null, loading: false });
    },
    setLoading: (loading: boolean) => {
      update(s => ({ ...s, loading }));
    },
    isAuthenticated: () => {
      if (typeof window === 'undefined') return false;
      return !!localStorage.getItem('token');
    }
  };
}

export const auth = createAuthStore();
