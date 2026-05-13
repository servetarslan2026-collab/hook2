<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import type { DeliveryAttempt, PaginatedResponse } from '$lib/api/client';

  let deadLetters = $state<DeliveryAttempt[]>([]);
  let total = $state(0);
  let currentPage = $state(1);
  let loading = $state(true);
  let actionLoading = $state<string | null>(null);
  let message = $state<{ type: 'success' | 'error'; text: string } | null>(null);

  onMount(() => loadDeadLetters());

  async function loadDeadLetters(page = 1) {
    loading = true;
    try {
      const result = await api.get(`admin/dead-letters?page=${page}&per_page=20`).json<PaginatedResponse<DeliveryAttempt>>();
      deadLetters = result.data;
      total = result.total;
      currentPage = page;
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function retryDelivery(id: string) {
    actionLoading = id;
    try {
      await api.post(`deliveries/${id}/retry`);
      message = { type: 'success', text: 'Retry queued' };
      await loadDeadLetters(currentPage);
    } catch (e: any) {
      message = { type: 'error', text: e?.message || 'Failed to retry' };
    } finally {
      actionLoading = null;
      setTimeout(() => message = null, 3000);
    }
  }
</script>

<div class="p-8">
  <div class="flex items-center justify-between mb-8">
    <div>
      <h1 class="text-2xl font-bold">Dead Letter Queue</h1>
      <p class="mt-1" style="color: var(--text-secondary);">Failed webhook deliveries that exhausted all retry attempts</p>
    </div>
    <div class="flex items-center gap-2">
      <span class="px-3 py-1.5 rounded-full text-sm font-medium"
        style="background: {total > 0 ? 'var(--danger)' : 'var(--success)'}20; color: {total > 0 ? 'var(--danger)' : 'var(--success)'};">
        {total} dead letter{total !== 1 ? 's' : ''}
      </span>
    </div>
  </div>

  {#if message}
    <div class="mb-4 px-4 py-3 rounded-lg text-sm"
      style="background: {message.type === 'success' ? 'var(--success)' : 'var(--danger)'}20; color: {message.type === 'success' ? 'var(--success)' : 'var(--danger)'};">
      {message.text}
    </div>
  {/if}

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if deadLetters.length === 0}
    <div class="text-center py-20 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <p class="text-4xl mb-3">🎉</p>
      <p class="mb-2 font-medium">No dead letters</p>
      <p style="color: var(--text-secondary);">All webhook deliveries are healthy</p>
    </div>
  {:else}
    <div class="rounded-xl border overflow-hidden" style="background: var(--bg-secondary); border-color: var(--border);">
      <table class="w-full">
        <thead>
          <tr style="background: var(--bg-tertiary);">
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Status</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Event ID</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Subscription ID</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">HTTP Code</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Attempts</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Time</th>
            <th class="px-4 py-3 text-right text-xs font-medium uppercase" style="color: var(--text-muted);">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y" style="border-color: var(--border);">
          {#each deadLetters as dl}
            <tr class="hover:bg-white/5 transition-colors">
              <td class="px-4 py-3">
                <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium"
                  style="background: var(--danger)20; color: var(--danger);">
                  <span class="w-1.5 h-1.5 rounded-full" style="background: var(--danger);"></span>
                  dead_letter
                </span>
              </td>
              <td class="px-4 py-3 text-xs font-mono" style="color: var(--text-muted);">{dl.event_id.slice(0, 8)}...</td>
              <td class="px-4 py-3 text-xs font-mono" style="color: var(--text-muted);">{dl.subscription_id.slice(0, 8)}...</td>
              <td class="px-4 py-3 text-sm font-mono">{dl.status_code || '-'}</td>
              <td class="px-4 py-3 text-sm">{dl.attempt_number}</td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-secondary);">
                {new Date(dl.created_at).toLocaleString()}
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex items-center justify-end gap-2">
                  <a href="/app/deliveries/{dl.id}" class="text-xs font-medium" style="color: var(--accent);">Details</a>
                  <button
                    onclick={() => retryDelivery(dl.id)}
                    disabled={actionLoading === dl.id}
                    class="px-3 py-1 rounded text-xs font-medium border text-white"
                    style="background: var(--accent); border-color: var(--accent);">
                    Retry
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
