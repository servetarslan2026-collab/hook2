<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { currentAppId } from '$lib/stores/app';
  import type { DeliveryAttempt, PaginatedResponse } from '$lib/api/client';

  let deliveries = $state<DeliveryAttempt[]>([]);
  let total = $state(0);
  let currentPage = $state(1);
  let statusFilter = $state('');
  let loading = $state(true);
  let appId = $state<string | null>(null);

  currentAppId.subscribe(id => { appId = id; });

  onMount(async () => {
    if (!appId) { goto('/'); return; }
    await loadDeliveries();
  });

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
</script>

<div class="p-8">
  <div class="flex items-center justify-between mb-8">
    <div>
      <h1 class="text-2xl font-bold">Deliveries</h1>
      <p class="mt-1" style="color: var(--text-secondary);">Webhook delivery logs and status</p>
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
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if deliveries.length === 0}
    <div class="text-center py-20 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <p class="mb-2 font-medium">No deliveries yet</p>
      <p style="color: var(--text-secondary);">Deliveries will appear here when events are sent</p>
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
          {#each deliveries as d}
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
