import { writable } from 'svelte store';

function createAppStore() {
  const stored = typeof window !== 'undefined' ? localStorage.getItem('currentAppId') : null;
  const { subscribe, set } = writable<string | null>(stored);

  return {
    subscribe,
    set: (id: string | null) => {
      if (typeof window !== 'undefined') {
        if (id) {
          localStorage.setItem('currentAppId', id);
        } else {
          localStorage.removeItem('currentAppId');
        }
      }
      set(id);
    },
    get: (): string | null => {
      if (typeof window !== 'undefined') {
        return localStorage.getItem('currentAppId');
      }
      return null;
    }
  };
}

export const currentAppId = createAppStore();
