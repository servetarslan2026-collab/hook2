<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import type { Organization, Application, StatsResponse } from '$lib/api/client';

  let orgs = $state<Organization[]>([]);
  let stats = $state<StatsResponse | null>(null);
  let loading = $state(true);

  onMount(async () => {
    if (!$auth.user) {
      goto('/login');
      return;
    }

    try {
      const { orgApi } = await import('$lib/api/client');
      orgs = await orgApi.list();

      if (orgs.length > 0) {
        stats = await orgApi.dashboard(orgs[0].id);
      }
    } catch (e) {
      console.error('Failed to load dashboard', e);
    } finally {
      loading = false;
    }
  });
</script>

<div class="p-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold">Dashboard</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Welcome back, {$auth.user?.name || 'User'}</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if orgs.length === 0}
    <!-- No organizations -->
    <div class="text-center py-20">
      <div class="w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4" style="background: var(--bg-tertiary);">
        <svg class="w-8 h-8" style="color: var(--text-muted);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
        </svg>
      </div>
      <h2 class="text-lg font-semibold mb-2">No organizations yet</h2>
      <p class="mb-6" style="color: var(--text-secondary);">Create your first organization to get started</p>
      <a href="/organizations/new" class="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium text-white"
        style="background: var(--accent);">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        Create Organization
      </a>
    </div>
  {:else}
    <!-- Stats -->
    {#if stats}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <p class="text-sm" style="color: var(--text-secondary);">Total Events</p>
          <p class="text-2xl font-bold mt-1">{stats.total_events.toLocaleString()}</p>
        </div>
        <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <p class="text-sm" style="color: var(--text-secondary);">Deliveries</p>
          <p class="text-2xl font-bold mt-1">{stats.total_deliveries.toLocaleString()}</p>
        </div>
        <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <p class="text-sm" style="color: var(--text-secondary);">Success Rate</p>
          <p class="text-2xl font-bold mt-1">{stats.success_rate.toFixed(1)}%</p>
        </div>
        <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <p class="text-sm" style="color: var(--text-secondary);">Subscriptions</p>
          <p class="text-2xl font-bold mt-1">{stats.total_subscriptions}</p>
        </div>
      </div>
    {/if}

    <!-- Organizations list -->
    <div class="mb-8">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold">Your Organizations</h2>
        <a href="/organizations/new" class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium text-white"
          style="background: var(--accent);">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          New
        </a>
      </div>

      <div class="space-y-3">
        {#each orgs as org}
          <a href="/organizations/{org.id}"
            class="flex items-center justify-between p-4 rounded-xl border transition-colors hover:border-gray-600"
            style="background: var(--bg-secondary); border-color: var(--border);">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center font-bold text-white" style="background: var(--accent);">
                {org.name[0].toUpperCase()}
              </div>
              <div>
                <p class="font-medium">{org.name}</p>
                <p class="text-sm" style="color: var(--text-muted);">Created {new Date(org.created_at).toLocaleDateString()}</p>
              </div>
            </div>
            <svg class="w-5 h-5" style="color: var(--text-muted);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
            </svg>
          </a>
        {/each}
      </div>
    </div>
  {/if}
</div>
