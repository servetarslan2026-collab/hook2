<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import type { Organization, PaginatedResponse } from '$lib/api/client';

  let orgs = $state<Organization[]>([]);
  let total = $state(0);
  let currentPage = $state(1);
  let loading = $state(true);
  let actionLoading = $state<string | null>(null);
  let message = $state<{ type: 'success' | 'error'; text: string } | null>(null);

  onMount(() => loadOrgs());

  async function loadOrgs(page = 1) {
    loading = true;
    try {
      const result = await api.get(`admin/organizations?page=${page}&per_page=20`).json<PaginatedResponse<Organization>>();
      orgs = result.data;
      total = result.total;
      currentPage = page;
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteOrg(org: Organization) {
    if (!confirm(`Delete organization "${org.name}"? All applications, events, and subscriptions will be deleted.`)) return;
    actionLoading = org.id;
    try {
      await api.delete(`admin/organizations/${org.id}`);
      message = { type: 'success', text: `Organization "${org.name}" deleted` };
      await loadOrgs(currentPage);
    } catch (e: any) {
      message = { type: 'error', text: e?.message || 'Failed to delete organization' };
    } finally {
      actionLoading = null;
      setTimeout(() => message = null, 3000);
    }
  }
</script>

<div class="p-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold">Organization Management</h1>
    <p class="mt-1" style="color: var(--text-secondary);">View and manage all organizations</p>
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
  {:else}
    <div class="rounded-xl border overflow-hidden" style="background: var(--bg-secondary); border-color: var(--border);">
      <table class="w-full">
        <thead>
          <tr style="background: var(--bg-tertiary);">
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Name</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">ID</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Owner</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Created</th>
            <th class="px-4 py-3 text-right text-xs font-medium uppercase" style="color: var(--text-muted);">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y" style="border-color: var(--border);">
          {#each orgs as org}
            <tr class="hover:bg-white/5 transition-colors">
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold text-white" style="background: var(--accent);">
                    {org.name[0]}
                  </div>
                  <span class="text-sm font-medium">{org.name}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-xs font-mono" style="color: var(--text-muted);">{org.id.slice(0, 8)}...</td>
              <td class="px-4 py-3 text-xs font-mono" style="color: var(--text-muted);">{org.owner_id.slice(0, 8)}...</td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-secondary);">
                {new Date(org.created_at).toLocaleDateString()}
              </td>
              <td class="px-4 py-3 text-right">
                <button
                  onclick={() => deleteOrg(org)}
                  disabled={actionLoading === org.id}
                  class="px-3 py-1 rounded text-xs font-medium border transition-colors"
                  style="border-color: var(--danger)40; color: var(--danger);">
                  Delete
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
