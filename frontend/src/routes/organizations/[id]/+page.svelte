<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import type { Application, StatsResponse, ChartDataPoint } from '$lib/api/client';

  let apps = $state<Application[]>([]);
  let stats = $state<(StatsResponse & { chart?: ChartDataPoint[] }) | null>(null);
  let loading = $state(true);
  let orgId = $derived($page.params.id);

  onMount(async () => {
    try {
      const { appApi, orgApi } = await import('$lib/api/client');
      apps = await appApi.list(orgId);
      stats = await orgApi.dashboard(orgId);
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });
</script>

<div class="p-8">
  <div class="flex items-center justify-between mb-8">
    <div>
      <h1 class="text-2xl font-bold">Organization</h1>
      <p class="mt-1" style="color: var(--text-secondary);">Manage your applications and settings</p>
    </div>
    <div class="flex gap-3">
      <a href="/organizations/{orgId}/settings" class="px-3 py-1.5 rounded-lg text-sm border" style="color: var(--text-secondary); border-color: var(--border);">
        Settings
      </a>
      <a href="/organizations/{orgId}/members" class="px-3 py-1.5 rounded-lg text-sm border" style="color: var(--text-secondary); border-color: var(--border);">
        Members
      </a>
    </div>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else}
    <!-- Stats -->
    {#if stats}
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <p class="text-sm" style="color: var(--text-secondary);">Events</p>
          <p class="text-2xl font-bold mt-1">{stats.total_events}</p>
        </div>
        <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <p class="text-sm" style="color: var(--text-secondary);">Deliveries</p>
          <p class="text-2xl font-bold mt-1">{stats.total_deliveries}</p>
        </div>
        <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <p class="text-sm" style="color: var(--text-secondary);">Success Rate</p>
          <p class="text-2xl font-bold mt-1">{stats.success_rate.toFixed(1)}%</p>
        </div>
        <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <p class="text-sm" style="color: var(--text-secondary);">Applications</p>
          <p class="text-2xl font-bold mt-1">{apps.length}</p>
        </div>
      </div>
    {/if}

    <!-- Applications -->
    <div>
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold">Applications</h2>
        <a href="/organizations/{orgId}/applications/new"
          class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium text-white"
          style="background: var(--accent);">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          New Application
        </a>
      </div>

      {#if apps.length === 0}
        <div class="text-center py-16 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <p class="mb-4" style="color: var(--text-secondary);">No applications yet</p>
          <a href="/organizations/{orgId}/applications/new" class="text-sm font-medium" style="color: var(--accent);">
            Create your first application →
          </a>
        </div>
      {:else}
        <div class="space-y-3">
          {#each apps as app}
            <a href="/app/{app.id}"
              class="flex items-center justify-between p-4 rounded-xl border transition-colors hover:border-gray-600"
              style="background: var(--bg-secondary); border-color: var(--border);">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg flex items-center justify-center font-bold text-white" style="background: var(--success);">
                  {app.name[0].toUpperCase()}
                </div>
                <div>
                  <p class="font-medium">{app.name}</p>
                  <p class="text-sm" style="color: var(--text-muted);">{app.description || 'No description'}</p>
                </div>
              </div>
              <svg class="w-5 h-5" style="color: var(--text-muted);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </a>
          {/each}
        </div>
      {/if}
    </div>
  {/if}
</div>
