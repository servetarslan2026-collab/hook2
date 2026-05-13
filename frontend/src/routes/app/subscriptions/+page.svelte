<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { currentAppId } from '$lib/stores/app';
  import type { Subscription } from '$lib/api/client';

  let subscriptions = $state<Subscription[]>([]);
  let loading = $state(true);
  let appId = $state<string | null>(null);

  currentAppId.subscribe(id => { appId = id; });

  onMount(async () => {
    if (!appId) { goto('/'); return; }
    try {
      const { subscriptionApi } = await import('$lib/api/client');
      subscriptions = await subscriptionApi.list(appId);
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  async function toggleSubscription(sub: Subscription) {
    try {
      const { subscriptionApi } = await import('$lib/api/client');
      await subscriptionApi.update(sub.id, { enabled: !sub.enabled });
      sub.enabled = !sub.enabled;
    } catch (e) {
      console.error(e);
    }
  }

  async function deleteSubscription(id: string) {
    if (!confirm('Delete this subscription?')) return;
    try {
      const { subscriptionApi } = await import('$lib/api/client');
      await subscriptionApi.delete(id);
      subscriptions = subscriptions.filter(s => s.id !== id);
    } catch (e) {
      console.error(e);
    }
  }
</script>

<div class="p-8">
  <div class="flex items-center justify-between mb-8">
    <div>
      <h1 class="text-2xl font-bold">Subscriptions</h1>
      <p class="mt-1" style="color: var(--text-secondary);">Manage webhook endpoints</p>
    </div>
    <a href="/app/subscriptions/new"
      class="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium text-white"
      style="background: var(--accent);">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
      </svg>
      New Subscription
    </a>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if subscriptions.length === 0}
    <div class="text-center py-20 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <p class="mb-2 font-medium">No subscriptions</p>
      <p class="mb-4" style="color: var(--text-secondary);">Create a subscription to start receiving webhooks</p>
      <a href="/app/subscriptions/new" class="text-sm font-medium" style="color: var(--accent);">Create subscription →</a>
    </div>
  {:else}
    <div class="space-y-3">
      {#each subscriptions as sub}
        <div class="p-4 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-3 h-3 rounded-full" style="background: {sub.enabled ? 'var(--success)' : 'var(--text-muted)'};"></div>
              <div>
                <p class="font-medium">{sub.description || 'Untitled Subscription'}</p>
                <p class="text-sm font-mono" style="color: var(--text-muted);">{sub.target_url}</p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <button onclick={() => toggleSubscription(sub)}
                class="px-3 py-1 rounded text-xs font-medium border"
                style="border-color: var(--border); color: {sub.enabled ? 'var(--success)' : 'var(--text-muted)'};">
                {sub.enabled ? 'Enabled' : 'Disabled'}
              </button>
              <button onclick={() => deleteSubscription(sub.id)} class="p-1.5 rounded hover:bg-red-500/10" title="Delete">
                <svg class="w-4 h-4" style="color: var(--danger);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                </svg>
              </button>
            </div>
          </div>
          <div class="mt-3 flex flex-wrap gap-2">
            {#each sub.event_types as et}
              <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium" style="background: rgba(59,130,246,0.1); color: var(--accent);">
                {et}
              </span>
            {/each}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
