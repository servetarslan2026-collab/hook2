<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import type { Application, StatsResponse, ChartDataPoint } from '$lib/api/client';

  let app = $state<Application | null>(null);
  let stats = $state<(StatsResponse & { chart?: ChartDataPoint[] }) | null>(null);
  let loading = $state(true);
  let appId = $derived($page.params.id);

  onMount(async () => {
    try {
      const { appApi } = await import('$lib/api/client');
      [app, stats] = await Promise.all([
        appApi.get(appId),
        appApi.dashboard(appId)
      ]);
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });
</script>

<div class="p-8">
  <div class="mb-8">
    <p class="text-sm" style="color: var(--text-muted);">Application</p>
    <h1 class="text-2xl font-bold">{app?.name || 'Loading...'}</h1>
    {#if app?.description}
      <p class="mt-1" style="color: var(--text-secondary);">{app.description}</p>
    {/if}
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else}
    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <div class="flex items-center gap-2 mb-2">
          <svg class="w-4 h-4" style="color: var(--accent);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
          </svg>
          <p class="text-sm" style="color: var(--text-secondary);">Total Events</p>
        </div>
        <p class="text-3xl font-bold">{stats?.total_events || 0}</p>
      </div>
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <div class="flex items-center gap-2 mb-2">
          <svg class="w-4 h-4" style="color: var(--success);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"/>
          </svg>
          <p class="text-sm" style="color: var(--text-secondary);">Deliveries</p>
        </div>
        <p class="text-3xl font-bold">{stats?.total_deliveries || 0}</p>
      </div>
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <div class="flex items-center gap-2 mb-2">
          <svg class="w-4 h-4" style="color: var(--success);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
          <p class="text-sm" style="color: var(--text-secondary);">Success Rate</p>
        </div>
        <p class="text-3xl font-bold">{stats?.success_rate?.toFixed(1) || 0}%</p>
      </div>
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <div class="flex items-center gap-2 mb-2">
          <svg class="w-4 h-4" style="color: var(--warning);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
          </svg>
          <p class="text-sm" style="color: var(--text-secondary);">Subscriptions</p>
        </div>
        <p class="text-3xl font-bold">{stats?.total_subscriptions || 0}</p>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
      <a href="/app/events/send"
        class="flex items-center gap-3 p-4 rounded-xl border transition-colors hover:border-gray-600"
        style="background: var(--bg-secondary); border-color: var(--border);">
        <div class="w-10 h-10 rounded-lg flex items-center justify-center" style="background: rgba(59,130,246,0.1);">
          <svg class="w-5 h-5" style="color: var(--accent);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
          </svg>
        </div>
        <div>
          <p class="font-medium">Send Test Event</p>
          <p class="text-sm" style="color: var(--text-muted);">Test your webhook setup</p>
        </div>
      </a>
      <a href="/app/subscriptions"
        class="flex items-center gap-3 p-4 rounded-xl border transition-colors hover:border-gray-600"
        style="background: var(--bg-secondary); border-color: var(--border);">
        <div class="w-10 h-10 rounded-lg flex items-center justify-center" style="background: rgba(34,197,94,0.1);">
          <svg class="w-5 h-5" style="color: var(--success);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
          </svg>
        </div>
        <div>
          <p class="font-medium">Manage Subscriptions</p>
          <p class="text-sm" style="color: var(--text-muted);">Configure webhook endpoints</p>
        </div>
      </a>
      <a href="/app/secrets"
        class="flex items-center gap-3 p-4 rounded-xl border transition-colors hover:border-gray-600"
        style="background: var(--bg-secondary); border-color: var(--border);">
        <div class="w-10 h-10 rounded-lg flex items-center justify-center" style="background: rgba(245,158,11,0.1);">
          <svg class="w-5 h-5" style="color: var(--warning);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
          </svg>
        </div>
        <div>
          <p class="font-medium">API Keys</p>
          <p class="text-sm" style="color: var(--text-muted);">Manage application secrets</p>
        </div>
      </a>
    </div>

    <!-- Recent Activity placeholder -->
    <div class="rounded-xl border p-6" style="background: var(--bg-secondary); border-color: var(--border);">
      <h2 class="text-lg font-semibold mb-4">Recent Activity</h2>
      <div class="text-center py-8">
        <p style="color: var(--text-muted);">No recent events. Send your first webhook to see activity here.</p>
      </div>
    </div>
  {/if}
</div>
