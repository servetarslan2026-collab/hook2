<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import { currentAppId } from '$lib/stores/app';
  import { ws, type DeliveryUpdate } from '$lib/stores/websocket';
  import type { DeliveryAttempt, PaginatedResponse } from '$lib/api/client';

  let deliveries = $state<DeliveryAttempt[]>([]);
  let total = $state(0);
  let currentPage = $state(1);
  let statusFilter = $state('');
  let loading = $state(true);
  let appId = $state<string | null>(null);
  let connected = $state(false);
  let newCount = $state(0);
  let toastQueue = $state<DeliveryUpdate[]>([]);

  let unsubWs: (() => void) | null = null;
  let unsubConnected: (() => void) | null = null;
  let toastTimer: ReturnType<typeof setTimeout> | null = null;

  currentAppId.subscribe(id => { appId = id; });

  onMount(async () => {
    if (!appId) { goto('/'); return; }

    // Subscribe to WebSocket connection status
    unsubConnected = ws.connected.subscribe(val => { connected = val; });

    // Connect WebSocket scoped to this app
    ws.subscribeApp(appId);

    // Listen for real-time delivery updates
    unsubWs = ws.onUpdate((update) => {
      // Only count if filtering matches
      if (!statusFilter || update.status === statusFilter) {
        newCount++;
      }

      // Show toast notification
      showToast(update);

      // Auto-refresh if on page 1 and no filter
      if (currentPage === 1 && !statusFilter) {
        prependUpdate(update);
      }
    });

    await loadDeliveries();
  });

  onDestroy(() => {
    unsubWs?.();
    unsubConnected?.();
    if (toastTimer) clearTimeout(toastTimer);
  });

  function showToast(update: DeliveryUpdate) {
    toastQueue = [...toastQueue, update].slice(-3);
    if (toastTimer) clearTimeout(toastTimer);
    toastTimer = setTimeout(() => { toastQueue = []; }, 4000);
  }

  function prependUpdate(update: DeliveryUpdate) {
    const attempt: DeliveryAttempt = {
      id: update.id,
      event_id: update.event_id,
      subscription_id: update.subscription_id,
      status: update.status,
      status_code: update.status_code,
      request_body: '',
      response_body: '',
      duration_ms: update.duration_ms,
      attempt_number: update.attempt_number,
      created_at: update.created_at,
    };
    deliveries = [attempt, ...deliveries].slice(0, 20);
    total++;
  }

  async function loadDeliveries(page = 1) {
    loading = true;
    try {
      const { deliveryApi } = await import('$lib/api/client');
      const params: any = { page, per_page: 20 };
      if (statusFilter) params.status = statusFilter;
      const result = await deliveryApi.list(appId, params);
      deliveries = result.data;
      total = result.total;
      currentPage = page;
      newCount = 0;
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  function statusColor(status: string) {
    switch (status) {
      case 'success': return 'var(--success)';
      case 'failed': return 'var(--danger)';
      case 'dead_letter': return 'var(--warning)';
      default: return 'var(--text-muted)';
    }
  }

  function statusIcon(status: string) {
    switch (status) {
      case 'success': return '✓';
      case 'failed': return '✗';
      case 'dead_letter': return '⚠';
      default: return '…';
    }
  }
</script>

<div class="p-8">
  <!-- Toast notifications -->
  {#if toastQueue.length > 0}
    <div class="fixed top-4 right-4 z-50 space-y-2">
      {#each toastQueue as toast (toast.id)}
        <div class="flex items-center gap-3 px-4 py-3 rounded-lg shadow-lg border text-sm animate-slide-in"
          style="background: var(--bg-secondary); border-color: {statusColor(toast.status)};">
          <span class="w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold text-white"
            style="background: {statusColor(toast.status)};">
            {statusIcon(toast.status)}
          </span>
          <div>
            <p class="font-medium">Webhook {toast.status}</p>
            <p class="text-xs" style="color: var(--text-muted);">HTTP {toast.status_code || '—'} · {toast.duration_ms}ms</p>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <div class="flex items-center justify-between mb-8">
    <div>
      <div class="flex items-center gap-3">
        <h1 class="text-2xl font-bold">Deliveries</h1>
        <!-- Live indicator -->
        <div class="flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium"
          style="background: {connected ? 'var(--success)' : 'var(--danger)'}20; color: {connected ? 'var(--success)' : 'var(--danger)'};">
          <span class="w-2 h-2 rounded-full {connected ? 'animate-pulse' : ''}"
            style="background: {connected ? 'var(--success)' : 'var(--danger)'};"></span>
          {connected ? 'LIVE' : 'OFFLINE'}
        </div>
        {#if newCount > 0 && (currentPage !== 1 || statusFilter)}
          <button onclick={() => loadDeliveries(1)}
            class="flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium text-white"
            style="background: var(--accent);">
            {newCount} new
          </button>
        {/if}
      </div>
      <p class="mt-1" style="color: var(--text-secondary);">Webhook delivery logs — real-time updates via WebSocket</p>
    </div>
    <div class="flex gap-2">
      <button onclick={() => { statusFilter = ''; loadDeliveries(); }}
        class="px-3 py-1.5 rounded-lg text-sm border {!statusFilter ? 'text-white' : ''}"
        style="background: {!statusFilter ? 'var(--accent)' : 'transparent'}; border-color: var(--border);">
        All
      </button>
      <button onclick={() => { statusFilter = 'success'; loadDeliveries(); }}
        class="px-3 py-1.5 rounded-lg text-sm border {statusFilter === 'success' ? 'text-white' : ''}"
        style="background: {statusFilter === 'success' ? 'var(--success)' : 'transparent'}; border-color: var(--border);">
        Success
      </button>
      <button onclick={() => { statusFilter = 'failed'; loadDeliveries(); }}
        class="px-3 py-1.5 rounded-lg text-sm border {statusFilter === 'failed' ? 'text-white' : ''}"
        style="background: {statusFilter === 'failed' ? 'var(--danger)' : 'transparent'}; border-color: var(--border);">
        Failed
      </button>
      <button onclick={() => { statusFilter = 'dead_letter'; loadDeliveries(); }}
        class="px-3 py-1.5 rounded-lg text-sm border {statusFilter === 'dead_letter' ? 'text-white' : ''}"
        style="background: {statusFilter === 'dead_letter' ? 'var(--warning)' : 'transparent'}; border-color: var(--border);">
        Dead Letter
      </button>
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if deliveries.length === 0}
    <div class="text-center py-20 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <p class="mb-2 font-medium">No deliveries yet</p>
      <p style="color: var(--text-secondary);">Deliveries will appear here in real-time when events are sent</p>
    </div>
  {:else}
    <div class="rounded-xl border overflow-hidden" style="background: var(--bg-secondary); border-color: var(--border);">
      <table class="w-full">
        <thead>
          <tr style="background: var(--bg-tertiary);">
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Status</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">HTTP Code</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Duration</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Attempt</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Time</th>
            <th class="px-4 py-3 text-right text-xs font-medium uppercase" style="color: var(--text-muted);">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y" style="border-color: var(--border);">
          {#each deliveries as d (d.id)}
            <tr class="hover:bg-white/5 transition-colors">
              <td class="px-4 py-3">
                <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium"
                  style="background: {statusColor(d.status)}20; color: {statusColor(d.status)};">
                  <span class="w-1.5 h-1.5 rounded-full" style="background: {statusColor(d.status)};"></span>
                  {d.status}
                </span>
              </td>
              <td class="px-4 py-3 text-sm font-mono">{d.status_code || '-'}</td>
              <td class="px-4 py-3 text-sm">{d.duration_ms}ms</td>
              <td class="px-4 py-3 text-sm">#{d.attempt_number}</td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-secondary);">
                {new Date(d.created_at).toLocaleString()}
              </td>
              <td class="px-4 py-3 text-right">
                <a href="/app/deliveries/{d.id}" class="text-sm font-medium" style="color: var(--accent);">Details →</a>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  @keyframes slide-in {
    from { transform: translateX(100%); opacity: 0; }
    to { transform: translateX(0); opacity: 1; }
  }
  .animate-slide-in {
    animation: slide-in 0.3s ease-out;
  }
</style>
