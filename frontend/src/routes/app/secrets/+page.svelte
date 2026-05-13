<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { currentAppId } from '$lib/stores/app';
  import type { ApplicationSecret } from '$lib/api/client';

  let secrets = $state<ApplicationSecret[]>([]);
  let loading = $state(true);
  let newName = $state('');
  let newSecret = $state<ApplicationSecret | null>(null);
  let appId = $state<string | null>(null);

  currentAppId.subscribe(id => { appId = id; });

  onMount(async () => {
    if (!appId) { goto('/'); return; }
    await loadSecrets();
  });

  async function loadSecrets() {
    try {
      const { secretApi } = await import('$lib/api/client');
      secrets = await secretApi.list(appId);
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function createSecret() {
    if (!newName.trim()) return;
    try {
      const { secretApi } = await import('$lib/api/client');
      newSecret = await secretApi.create(appId, { name: newName });
      secrets = [newSecret, ...secrets];
      newName = '';
    } catch (e) {
      console.error(e);
    }
  }

  async function deleteSecret(id: string) {
    if (!confirm('Delete this secret? This cannot be undone.')) return;
    try {
      const { secretApi } = await import('$lib/api/client');
      await secretApi.delete(appId, id);
      secrets = secrets.filter(s => s.id !== id);
    } catch (e) {
      console.error(e);
    }
  }

  function copyKey(key: string) {
    navigator.clipboard.writeText(key);
  }
</script>

<div class="p-8">
  <div class="mb-8">
    <h1 class="text-2xl font-bold">Application Secrets</h1>
    <p class="mt-1" style="color: var(--text-secondary);">API keys for authenticating webhook requests</p>
  </div>

  <!-- Create new -->
  <div class="p-4 rounded-xl border mb-6" style="background: var(--bg-secondary); border-color: var(--border);">
    <h3 class="text-sm font-medium mb-3">Create New Secret</h3>
    <form onsubmit={(e) => { e.preventDefault(); createSecret(); }} class="flex gap-3">
      <input type="text" bind:value={newName} placeholder="Secret name (e.g., Production)"
        class="flex-1 px-4 py-2 rounded-lg border text-sm focus:outline-none focus:ring-2"
        style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary);" />
      <button type="submit" class="px-4 py-2 rounded-lg text-sm font-medium text-white" style="background: var(--accent);">
        Create
      </button>
    </form>
  </div>

  <!-- Newly created secret banner -->
  {#if newSecret}
    <div class="p-4 rounded-xl border mb-6" style="background: rgba(34,197,94,0.1); border-color: rgba(34,197,94,0.3);">
      <p class="text-sm font-medium mb-2" style="color: var(--success);">✓ Secret created! Copy it now — it won't be shown again.</p>
      <div class="flex items-center gap-2 p-2 rounded" style="background: var(--bg-tertiary);">
        <code class="flex-1 text-sm font-mono break-all">{newSecret.key}</code>
        <button onclick={() => copyKey(newSecret!.key)} class="px-3 py-1 rounded text-xs font-medium border" style="border-color: var(--border);">
          Copy
        </button>
      </div>
    </div>
  {/if}

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if secrets.length === 0}
    <div class="text-center py-16 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <p class="mb-2 font-medium">No secrets</p>
      <p style="color: var(--text-secondary);">Create an API key to authenticate your requests</p>
    </div>
  {:else}
    <div class="space-y-3">
      {#each secrets as secret}
        <div class="flex items-center justify-between p-4 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <div>
            <p class="font-medium">{secret.name}</p>
            <p class="text-sm font-mono" style="color: var(--text-muted);">{secret.key.slice(0, 12)}...{secret.key.slice(-4)}</p>
            <p class="text-xs mt-1" style="color: var(--text-muted);">Created {new Date(secret.created_at).toLocaleDateString()}</p>
          </div>
          <button onclick={() => deleteSecret(secret.id)} class="p-2 rounded hover:bg-red-500/10" title="Delete">
            <svg class="w-4 h-4" style="color: var(--danger);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>
