<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import type { User, PaginatedResponse } from '$lib/api/client';

  let users = $state<User[]>([]);
  let total = $state(0);
  let currentPage = $state(1);
  let loading = $state(true);
  let actionLoading = $state<string | null>(null);
  let message = $state<{ type: 'success' | 'error'; text: string } | null>(null);

  onMount(() => loadUsers());

  async function loadUsers(page = 1) {
    loading = true;
    try {
      const result = await api.get(`admin/users?page=${page}&per_page=20`).json<PaginatedResponse<User>>();
      users = result.data;
      total = result.total;
      currentPage = page;
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function toggleAdmin(user: User) {
    actionLoading = user.id;
    try {
      await api.put(`admin/users/${user.id}/admin`, { json: { is_admin: !user.is_admin } });
      message = { type: 'success', text: `${user.name} is now ${user.is_admin ? 'a regular user' : 'an admin'}` };
      await loadUsers(currentPage);
    } catch (e: any) {
      message = { type: 'error', text: e?.message || 'Failed to update user' };
    } finally {
      actionLoading = null;
      setTimeout(() => message = null, 3000);
    }
  }

  async function deleteUser(user: User) {
    if (!confirm(`Delete user "${user.name}" (${user.email})? This cannot be undone.`)) return;
    actionLoading = user.id;
    try {
      await api.delete(`admin/users/${user.id}`);
      message = { type: 'success', text: `User ${user.name} deleted` };
      await loadUsers(currentPage);
    } catch (e: any) {
      message = { type: 'error', text: e?.message || 'Failed to delete user' };
    } finally {
      actionLoading = null;
      setTimeout(() => message = null, 3000);
    }
  }
</script>

<div class="p-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold">User Management</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Manage all registered users</p>
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
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">User</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Email</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Role</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Verified</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Joined</th>
            <th class="px-4 py-3 text-right text-xs font-medium uppercase" style="color: var(--text-muted);">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y" style="border-color: var(--border);">
          {#each users as user}
            <tr class="hover:bg-white/5 transition-colors">
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium text-white"
                    style="background: {user.is_admin ? 'var(--danger)' : 'var(--accent)'};">
                    {user.name?.[0] || 'U'}
                  </div>
                  <span class="text-sm font-medium">{user.name}</span>
                </div>
              </td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-secondary);">{user.email}</td>
              <td class="px-4 py-3">
                {#if user.is_admin}
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    style="background: var(--danger)20; color: var(--danger);">
                    Admin
                  </span>
                {:else}
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    style="background: var(--accent)20; color: var(--accent);">
                    User
                  </span>
                {/if}
              </td>
              <td class="px-4 py-3">
                {#if user.email_verified}
                  <span style="color: var(--success);">✓</span>
                {:else}
                  <span style="color: var(--text-muted);">✗</span>
                {/if}
              </td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-secondary);">
                {new Date(user.created_at).toLocaleDateString()}
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button
                    onclick={() => toggleAdmin(user)}
                    disabled={actionLoading === user.id}
                    class="px-3 py-1 rounded text-xs font-medium border transition-colors"
                    style="border-color: var(--border); color: var(--text-secondary);">
                    {user.is_admin ? 'Remove Admin' : 'Make Admin'}
                  </button>
                  <button
                    onclick={() => deleteUser(user)}
                    disabled={actionLoading === user.id}
                    class="px-3 py-1 rounded text-xs font-medium border transition-colors"
                    style="border-color: var(--danger)40; color: var(--danger);">
                    Delete
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
