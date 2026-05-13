<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import type { DeliveryAttempt } from '$lib/api/client';

  let delivery = $state<DeliveryAttempt | null>(null);
  let loading = $state(true);
  let retrying = $state(false);
  let deliveryId = $derived($page.params.id);

  onMount(async () => {
    try {
      const { deliveryApi } = await import('$lib/api/client');
      delivery = await deliveryApi.get(deliveryId);
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  async function handleRetry() {
    retrying = true;
    try {
      const { deliveryApi } = await import('$lib/api/client');
      await deliveryApi.retry(deliveryId);
      delivery = await deliveryApi.get(deliveryId);
    } catch (e) {
      console.error(e);
    } finally {
      retrying = false;
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

  function formatJson(str: string) {
    try {
      return JSON.stringify(JSON.parse(str), null, 2);
    } catch {
      return str || '(empty)';
    }
  }
</script>

<div class="p-8 max-w-4xl">
  <div class="mb-8">
    <div class="flex items-center gap-2 mb-2">
      <a href="/app/deliveries" class="text-sm" style="color: var(--text-muted);">← Back to deliveries</a>
    </div>
    <h1 class="text-2xl font-bold">Delivery Details</h1>
    <p class="mt-1 font-mono text-sm" style="color: var(--text-muted);">{deliveryId}</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if !delivery}
    <div class="text-center py-20 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <p class="text-lg font-medium">Delivery not found</p>
      <p class="mt-2" style="color: var(--text-secondary);">This delivery may have been cleaned up.</p>
    </div>
  {:else}
    <!-- Status Overview -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <p class="text-sm" style="color: var(--text-secondary);">Status</p>
        <div class="flex items-center gap-2 mt-2">
          <span class="w-2.5 h-2.5 rounded-full" style="background: {statusColor(delivery.status)};"></span>
          <span class="font-semibold capitalize">{delivery.status}</span>
        </div>
      </div>
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <p class="text-sm" style="color: var(--text-secondary);">HTTP Status</p>
        <p class="text-2xl font-bold mt-1">{delivery.status_code || '—'}</p>
      </div>
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <p class="text-sm" style="color: var(--text-secondary);">Duration</p>
        <p class="text-2xl font-bold mt-1">{delivery.duration_ms}ms</p>
      </div>
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <p class="text-sm" style="color: var(--text-secondary);">Attempt</p>
        <p class="text-2xl font-bold mt-1">#{delivery.attempt_number}</p>
      </div>
    </div>

    <!-- Meta -->
    <div class="p-4 rounded-xl border mb-6" style="background: var(--bg-secondary); border-color: var(--border);">
      <div class="grid grid-cols-2 gap-4">
        <div>
          <p class="text-xs uppercase" style="color: var(--text-muted);">Event ID</p>
          <p class="text-sm font-mono mt-1">{delivery.event_id}</p>
        </div>
        <div>
          <p class="text-xs uppercase" style="color: var(--text-muted);">Subscription ID</p>
          <p class="text-sm font-mono mt-1">{delivery.subscription_id}</p>
        </div>
        <div>
          <p class="text-xs uppercase" style="color: var(--text-muted);">Created</p>
          <p class="text-sm mt-1">{new Date(delivery.created_at).toLocaleString()}</p>
        </div>
        <div>
          <p class="text-xs uppercase" style="color: var(--text-muted);">Actions</p>
          <button onclick={handleRetry} disabled={retrying}
            class="mt-1 px-3 py-1 rounded text-sm font-medium text-white disabled:opacity-50"
            style="background: var(--accent);">
            {retrying ? 'Retrying...' : 'Retry Delivery'}
          </button>
        </div>
      </div>
    </div>

    <!-- Request Body -->
    <div class="mb-6">
      <h3 class="text-sm font-medium mb-2" style="color: var(--text-secondary);">Request Body (sent to webhook URL)</h3>
      <div class="rounded-xl border overflow-hidden" style="background: var(--bg-secondary); border-color: var(--border);">
        <pre class="p-4 text-sm font-mono overflow-x-auto" style="color: var(--text-primary);">{formatJson(delivery.request_body)}</pre>
      </div>
    </div>

    <!-- Response Body -->
    <div>
      <h3 class="text-sm font-medium mb-2" style="color: var(--text-secondary);">Response Body (from webhook URL)</h3>
      <div class="rounded-xl border overflow-hidden" style="background: var(--bg-secondary); border-color: var(--border);">
        <pre class="p-4 text-sm font-mono overflow-x-auto" style="color: var(--text-primary);">{formatJson(delivery.response_body)}</pre>
      </div>
    </div>
  {/if}
</div>
