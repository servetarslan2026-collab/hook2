<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import type { Organization } from '$lib/api/client';

  let org = $state<Organization | null>(null);
  let name = $state('');
  let loading = $state(true);
  let saving = $state(false);
  let error = $state('');
  let success = $state('');
  let orgId = $derived($page.params.id);

  onMount(async () => {
    try {
      const { orgApi } = await import('$lib/api/client');
      org = await orgApi.get(orgId);
      name = org.name;
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  async function handleUpdate() {
    error = '';
    success = '';
    saving = true;
    try {
      const { orgApi } = await import('$lib/api/client');
      org = await orgApi.update(orgId, { name });
      success = 'Organization updated successfully';
    } catch (e: any) {
      error = e.message || 'Failed to update organization';
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!confirm('Are you sure you want to delete this organization? This action cannot be undone.')) return;
    try {
      const { orgApi } = await import('$lib/api/client');
      await orgApi.delete(orgId);
      goto('/organizations');
    } catch (e: any) {
      error = e.message || 'Failed to delete organization';
    }
  }
</script>

<div class="p-8 max-w-2xl">
  <div class="mb-8">
    <div class="flex items-center gap-2 mb-2">
      <a href="/organizations/{orgId}" class="text-sm" style="color: var(--text-muted);">← Back</a>
    </div>
    <h1 class="text-2xl font-bold">Organization Settings</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Manage your organization details</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else}
    <!-- Edit Form -->
    <div class="p-6 rounded-xl border mb-6" style="background: var(--bg-secondary); border-color: var(--border);">
      <h2 class="text-lg font-semibold mb-4">General</h2>
      <form onsubmit={(e) => { e.preventDefault(); handleUpdate(); }} class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Organization Name</label>
          <input type="text" bind:value={name} placeholder="Organization name" required
            class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
            style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
        </div>

        {#if error}
          <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(239,68,68,0.1); color: var(--danger);">
            {error}
          </div>
        {/if}

        {#if success}
          <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(34,197,94,0.1); color: var(--success);">
            {success}
          </div>
        {/if}

        <button type="submit" disabled={saving}
          class="px-4 py-2 rounded-lg text-sm font-medium text-white disabled:opacity-50"
          style="background: var(--accent);">
          {saving ? 'Saving...' : 'Save Changes'}
        </button>
      </form>
    </div>

    <!-- Danger Zone -->
    <div class="p-6 rounded-xl border" style="background: var(--bg-secondary); border-color: rgba(239,68,68,0.3);">
      <h2 class="text-lg font-semibold mb-2" style="color: var(--danger);">Danger Zone</h2>
      <p class="text-sm mb-4" style="color: var(--text-secondary);">
        Deleting an organization will permanently remove all applications, subscriptions, events, and data associated with it.
      </p>
      <button onclick={handleDelete}
        class="px-4 py-2 rounded-lg text-sm font-medium text-white"
        style="background: var(--danger);">
        Delete Organization
      </button>
    </div>
  {/if}
</div>
