<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import type { Organization } from '$lib/api/client';

  let orgs = $state<Organization[]>([]);
  let loading = $state(true);

  onMount(async () => {
    if (!$auth.user) {
      goto('/login');
      return;
    }
    try {
      const { orgApi } = await import('$lib/api/client');
      orgs = await orgApi.list();
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
      <h1 class="text-2xl font-bold">Organizations</h1>
      <p class="mt-1" style="color: var(--text-secondary);">Manage your organizations and workspaces</p>
    </div>
    <a href="/organizations/new"
      class="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium text-white"
      style="background: var(--accent);">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
      </svg>
      New Organization
    </a>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if orgs.length === 0}
    <div class="text-center py-20 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
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
    <div class="space-y-3">
      {#each orgs as org}
        <a href="/organizations/{org.id}"
          class="flex items-center justify-between p-4 rounded-xl border transition-colors hover:border-gray-600"
          style="background: var(--bg-secondary); border-color: var(--border);">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-xl flex items-center justify-center font-bold text-white text-lg" style="background: var(--accent);">
              {org.name[0].toUpperCase()}
            </div>
            <div>
              <p class="font-medium text-lg">{org.name}</p>
              <p class="text-sm" style="color: var(--text-muted);">Created {new Date(org.created_at).toLocaleDateString()}</p>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <a href="/organizations/{org.id}/settings" class="px-3 py-1.5 rounded-lg text-xs border" style="color: var(--text-secondary); border-color: var(--border);" onclick={(e) => e.stopPropagation()}>
              Settings
            </a>
            <a href="/organizations/{org.id}/members" class="px-3 py-1.5 rounded-lg text-xs border" style="color: var(--text-secondary); border-color: var(--border);" onclick={(e) => e.stopPropagation()}>
              Members
            </a>
            <svg class="w-5 h-5" style="color: var(--text-muted);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
            </svg>
          </div>
        </a>
      {/each}
    </div>
  {/if}
</div>
