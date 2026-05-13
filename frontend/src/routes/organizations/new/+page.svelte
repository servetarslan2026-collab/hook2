<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import type { Organization } from '$lib/api/client';

  let name = $state('');
  let loading = $state(false);
  let error = $state('');

  async function handleSubmit() {
    error = '';
    loading = true;
    try {
      const { orgApi } = await import('$lib/api/client');
      const org = await orgApi.create({ name });
      goto(`/organizations/${org.id}`);
    } catch (e: any) {
      error = e.message || 'Failed to create organization';
    } finally {
      loading = false;
    }
  }
</script>

<div class="p-8 max-w-xl">
  <div class="mb-8">
    <h1 class="text-2xl font-bold">Create Organization</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Set up a new workspace for your webhooks</p>
  </div>

  <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
    <div>
      <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Organization Name</label>
      <input type="text" bind:value={name} placeholder="My Company" required
        class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
        style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
    </div>

    {#if error}
      <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(239,68,68,0.1); color: var(--danger);">
        {error}
      </div>
    {/if}

    <div class="flex gap-3">
      <button type="submit" disabled={loading}
        class="px-4 py-2 rounded-lg text-sm font-medium text-white disabled:opacity-50"
        style="background: var(--accent);">
        {loading ? 'Creating...' : 'Create Organization'}
      </button>
      <a href="/" class="px-4 py-2 rounded-lg text-sm font-medium border" style="color: var(--text-secondary); border-color: var(--border);">
        Cancel
      </a>
    </div>
  </form>
</div>
