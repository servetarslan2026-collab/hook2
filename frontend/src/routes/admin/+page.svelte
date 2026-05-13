<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';

  let stats = $state<any>(null);
  let loading = $state(true);

  onMount(async () => {
    try {
      stats = await api.get('admin/stats').json();
    } catch (e) {
      console.error('Failed to load admin stats:', e);
    } finally {
      loading = false;
    }
  });

  function statCard(label: string, value: string | number, color: string, icon: string) {
    return { label, value, color, icon };
  }

  let cards = $derived(stats ? [
    statCard('Users', stats.total_users, 'var(--accent)', '👤'),
    statCard('Organizations', stats.total_organizations, 'var(--accent)', '🏢'),
    statCard('Applications', stats.total_applications, 'var(--success)', '📱'),
    statCard('Events', stats.total_events, 'var(--success)', '⚡'),
    statCard('Deliveries', stats.total_deliveries, 'var(--warning)', '📦'),
    statCard('Subscriptions', stats.total_subscriptions, 'var(--warning)', '🔔'),
    statCard('Success Rate', `${stats.success_rate.toFixed(1)}%`, stats.success_rate >= 95 ? 'var(--success)' : 'var(--danger)', '✅'),
    statCard('Dead Letters', stats.dead_letter_count, stats.dead_letter_count > 0 ? 'var(--danger)' : 'var(--success)', '💀'),
  ] : []);
</script>

<div class="p-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold">System Overview</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Real-time system health and statistics</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if stats}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {#each cards as card}
        <div class="rounded-xl border p-6" style="background: var(--bg-secondary); border-color: var(--border);">
          <div class="flex items-center justify-between mb-3">
            <span class="text-2xl">{card.icon}</span>
            <span class="text-xs font-medium px-2 py-0.5 rounded-full"
              style="background: {card.color}20; color: {card.color};">
              {card.label}
            </span>
          </div>
          <p class="text-3xl font-bold">{card.value}</p>
        </div>
      {/each}
    </div>

    <!-- Health indicators -->
    <div class="mt-8 grid grid-cols-1 lg:grid-cols-3 gap-4">
      <div class="rounded-xl border p-6" style="background: var(--bg-secondary); border-color: var(--border);">
        <h3 class="font-semibold mb-4 flex items-center gap-2">
          <span class="w-3 h-3 rounded-full" style="background: {stats.success_rate >= 95 ? 'var(--success)' : 'var(--danger)'};"></span>
          Delivery Health
        </h3>
        <div class="space-y-3">
          <div class="flex justify-between text-sm">
            <span style="color: var(--text-secondary);">Success Rate</span>
            <span class="font-mono">{stats.success_rate.toFixed(1)}%</span>
          </div>
          <div class="w-full h-2 rounded-full" style="background: var(--border);">
            <div class="h-full rounded-full transition-all"
              style="width: {stats.success_rate}%; background: {stats.success_rate >= 95 ? 'var(--success)' : 'var(--danger)'};"></div>
          </div>
          <div class="flex justify-between text-sm">
            <span style="color: var(--text-secondary);">Dead Letters</span>
            <span class="font-mono" style="color: {stats.dead_letter_count > 0 ? 'var(--danger)' : 'var(--success)'};">
              {stats.dead_letter_count}
            </span>
          </div>
        </div>
      </div>

      <div class="rounded-xl border p-6" style="background: var(--bg-secondary); border-color: var(--border);">
        <h3 class="font-semibold mb-4">📊 Content Summary</h3>
        <div class="space-y-3 text-sm">
          <div class="flex justify-between">
            <span style="color: var(--text-secondary);">Users</span>
            <span class="font-mono">{stats.total_users}</span>
          </div>
          <div class="flex justify-between">
            <span style="color: var(--text-secondary);">Organizations</span>
            <span class="font-mono">{stats.total_organizations}</span>
          </div>
          <div class="flex justify-between">
            <span style="color: var(--text-secondary);">Applications</span>
            <span class="font-mono">{stats.total_applications}</span>
          </div>
          <div class="flex justify-between">
            <span style="color: var(--text-secondary);">Subscriptions</span>
            <span class="font-mono">{stats.total_subscriptions}</span>
          </div>
        </div>
      </div>

      <div class="rounded-xl border p-6" style="background: var(--bg-secondary); border-color: var(--border);">
        <h3 class="font-semibold mb-4">⚡ Activity</h3>
        <div class="space-y-3 text-sm">
          <div class="flex justify-between">
            <span style="color: var(--text-secondary);">Total Events</span>
            <span class="font-mono">{stats.total_events.toLocaleString()}</span>
          </div>
          <div class="flex justify-between">
            <span style="color: var(--text-secondary);">Total Deliveries</span>
            <span class="font-mono">{stats.total_deliveries.toLocaleString()}</span>
          </div>
          <div class="flex justify-between">
            <span style="color: var(--text-secondary);">Avg Deliveries/Event</span>
            <span class="font-mono">{stats.total_events > 0 ? (stats.total_deliveries / stats.total_events).toFixed(1) : '0'}</span>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>
