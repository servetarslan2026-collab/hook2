<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import type { Organization } from '$lib/api/client';

  let org = $state<Organization | null>(null);
  let members = $state<any[]>([]);
  let loading = $state(true);
  let inviteEmail = $state('');
  let inviteRole = $state('member');
  let inviting = $state(false);
  let error = $state('');
  let orgId = $derived($page.params.id);

  onMount(async () => {
    try {
      const { orgApi } = await import('$lib/api/client');
      [org, members] = await Promise.all([
        orgApi.get(orgId),
        orgApi.members(orgId)
      ]);
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  async function handleInvite() {
    if (!inviteEmail.trim()) return;
    error = '';
    inviting = true;
    try {
      const { orgApi } = await import('$lib/api/client');
      await orgApi.inviteMember(orgId, { email: inviteEmail, role: inviteRole });
      members = await orgApi.members(orgId);
      inviteEmail = '';
    } catch (e: any) {
      error = e.message || 'Failed to invite member';
    } finally {
      inviting = false;
    }
  }

  async function removeMember(userId: string) {
    if (!confirm('Remove this member from the organization?')) return;
    try {
      const { orgApi } = await import('$lib/api/client');
      await orgApi.removeMember(orgId, userId);
      members = members.filter(m => m.id !== userId);
    } catch (e: any) {
      error = e.message || 'Failed to remove member';
    }
  }
</script>

<div class="p-8 max-w-3xl">
  <div class="mb-8">
    <div class="flex items-center gap-2 mb-2">
      <a href="/organizations/{orgId}" class="text-sm" style="color: var(--text-muted);">← Back</a>
    </div>
    <h1 class="text-2xl font-bold">Members</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Manage who has access to {org?.name || 'this organization'}</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else}
    <!-- Invite Form -->
    <div class="p-4 rounded-xl border mb-6" style="background: var(--bg-secondary); border-color: var(--border);">
      <h3 class="text-sm font-medium mb-3">Invite Member</h3>
      <form onsubmit={(e) => { e.preventDefault(); handleInvite(); }} class="flex gap-3">
        <input type="email" bind:value={inviteEmail} placeholder="email@example.com" required
          class="flex-1 px-4 py-2 rounded-lg border text-sm focus:outline-none focus:ring-2"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary);" />
        <select bind:value={inviteRole}
          class="px-3 py-2 rounded-lg border text-sm focus:outline-none"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary);">
          <option value="member">Member</option>
          <option value="admin">Admin</option>
        </select>
        <button type="submit" disabled={inviting}
          class="px-4 py-2 rounded-lg text-sm font-medium text-white disabled:opacity-50"
          style="background: var(--accent);">
          {inviting ? 'Inviting...' : 'Invite'}
        </button>
      </form>
      {#if error}
        <p class="mt-2 text-sm" style="color: var(--danger);">{error}</p>
      {/if}
    </div>

    <!-- Members List -->
    {#if members.length === 0}
      <div class="text-center py-16 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <p style="color: var(--text-secondary);">No members yet. Invite someone to get started.</p>
      </div>
    {:else}
      <div class="space-y-3">
        {#each members as member}
          <div class="flex items-center justify-between p-4 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-medium text-white" style="background: var(--accent);">
                {(member.name || member.email || 'U')[0].toUpperCase()}
              </div>
              <div>
                <p class="font-medium">{member.name || 'Unknown'}</p>
                <p class="text-sm" style="color: var(--text-muted);">{member.email}</p>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <span class="px-2 py-0.5 rounded text-xs font-medium" style="background: var(--bg-tertiary); color: var(--text-secondary);">
                {member.role || 'member'}
              </span>
              <button onclick={() => removeMember(member.id)} class="p-1.5 rounded hover:bg-red-500/10" title="Remove">
                <svg class="w-4 h-4" style="color: var(--danger);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>
