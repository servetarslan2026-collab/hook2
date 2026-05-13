import { writable, get } from 'svelte/store';
import { PUBLIC_API_URL } from '$env/static/public';
import { auth } from './auth';

export interface DeliveryUpdate {
  type: string;
  id: string;
  event_id: string;
  subscription_id: string;
  status: string;
  status_code: number;
  duration_ms: number;
  attempt_number: number;
  created_at: string;
}

type WSHandler = (update: DeliveryUpdate) => void;

function createWebSocketStore() {
  const connected = writable(false);
  const updates = writable<DeliveryUpdate[]>([]);
  let socket: WebSocket | null = null;
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  let handlers: WSHandler[] = [];
  let currentAppId: string | null = null;

  function getWsUrl(appId?: string): string {
    const apiUrl = PUBLIC_API_URL || '/api/v1';
    // Convert http(s) to ws(s)
    const wsBase = apiUrl.replace(/^http/, 'ws');
    const base = typeof window !== 'undefined'
      ? `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${wsBase}`
      : wsBase;
    const url = `${base}/ws/deliveries`;
    return appId ? `${url}?app_id=${appId}` : url;
  }

  function connect(appId?: string) {
    if (socket && socket.readyState === WebSocket.OPEN) {
      // Already connected, maybe switch subscription
      if (appId && appId !== currentAppId) {
        currentAppId = appId;
        socket.send(JSON.stringify({ action: 'subscribe', app_id: appId }));
      }
      return;
    }

    // Clean up existing connection
    disconnect();
    currentAppId = appId || null;

    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    if (!token) return;

    const url = getWsUrl(appId);
    socket = new WebSocket(url);

    socket.onopen = () => {
      connected.set(true);
      if (reconnectTimer) {
        clearTimeout(reconnectTimer);
        reconnectTimer = null;
      }
    };

    socket.onmessage = (event) => {
      try {
        const update: DeliveryUpdate = JSON.parse(event.data);
        if (update.type === 'delivery_update') {
          // Add to updates list (keep last 100)
          updates.update(list => [update, ...list].slice(0, 100));
          // Notify handlers
          for (const handler of handlers) {
            handler(update);
          }
        }
      } catch (e) {
        console.error('Failed to parse WS message:', e);
      }
    };

    socket.onclose = () => {
      connected.set(false);
      socket = null;
      // Reconnect after 3 seconds
      reconnectTimer = setTimeout(() => connect(currentAppId || undefined), 3000);
    };

    socket.onerror = (err) => {
      console.error('WebSocket error:', err);
    };
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
    if (socket) {
      socket.close();
      socket = null;
    }
    connected.set(false);
  }

  function subscribeApp(appId: string) {
    currentAppId = appId;
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify({ action: 'subscribe', app_id: appId }));
    } else {
      connect(appId);
    }
  }

  function onUpdate(handler: WSHandler) {
    handlers.push(handler);
    return () => {
      handlers = handlers.filter(h => h !== handler);
    };
  }

  function clearUpdates() {
    updates.set([]);
  }

  return {
    connected,
    updates,
    connect,
    disconnect,
    subscribeApp,
    onUpdate,
    clearUpdates,
  };
}

export const ws = createWebSocketStore();
