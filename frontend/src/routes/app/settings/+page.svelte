<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { currentAppId } from '$lib/stores/app';
  import type { Application } from '$lib/api/client';

  let app = $state<Application | null>(null);
  let name = $state('');
  let description = $state('');
  let loading = $state(true);
  let saving = $state(false);
  let error = $state('');
  let success = $state('');
  let appId = $state<string | null>(null);

  currentAppId.subscribe(id => { appId = id; });

  onMount(async () => {
    if (!appId) { goto('/'); return; }
    try {
      const { appApi } = await import('$lib/api/client');
      app = await appApi.get(appId);
      name = app.name;
      description = app.description || '';
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  async function handleUpdate() {
    if (!appId) return;
    error = '';
    success = '';
    saving = true;
    try {
      const { appApi } = await import('$lib/api/client');
      app = await appApi.update(appId, { name, description });
      success = 'Application updated successfully';
    } catch (e: any) {
      error = e.message || 'Failed to update application';
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!appId) return;
    if (!confirm('Are you sure you want to delete this application? All events, subscriptions, and delivery logs will be permanently removed.')) return;
    try {
      const { appApi } = await import('$lib/api/client');
      await appApi.delete(appId);
      currentAppId.set(null);
      goto('/');
    } catch (e: any) {
      error = e.message || 'Failed to delete application';
    }
  }
</script>

<div class="p-8 max-w-2xl">
  <div class="mb-8">
    <h1 class="text-2xl font-bold">Application Settings</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Configure {app?.name || 'this application'}</p>
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
          <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Application Name</label>
          <input type="text" bind:value={name} placeholder="My App" required
            class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
            style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Description</label>
          <textarea bind:value={description} placeholder="What does this application do?" rows="3"
            class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2 resize-none"
            style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);"></textarea>
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
        Deleting this application will permanently remove all events, subscriptions, delivery logs, and API keys.
      </p>
      <button onclick={handleDelete}
        class="px-4 py-2 rounded-lg text-sm font-medium text-white"
        style="background: var(--danger);">
        Delete Application
      </button>
    </div>
  {/if}
</div>
