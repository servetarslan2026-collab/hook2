<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';

  let tokens = $state<any[]>([]);
  let loading = $state(true);
  let newName = $state('');
  let creating = $state(false);
  let newToken = $state<string | null>(null);
  let orgId = $derived($page.params.id);

  onMount(async () => {
    try {
      // Service tokens API not yet in client, placeholder
      tokens = [];
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  function copyToken(token: string) {
    navigator.clipboard.writeText(token);
  }
</script>

<div class="p-8 max-w-3xl">
  <div class="mb-8">
    <div class="flex items-center gap-2 mb-2">
      <a href="/organizations/{orgId}" class="text-sm" style="color: var(--text-muted);">← Back</a>
    </div>
    <h1 class="text-2xl font-bold">Service Tokens</h1>
    <p class="mt-1" style="color: var(--text-secondary);">API tokens for programmatic access to this organization</p>
  </div>

  <!-- Create Form -->
  <div class="p-4 rounded-xl border mb-6" style="background: var(--bg-secondary); border-color: var(--border);">
    <h3 class="text-sm font-medium mb-3">Create New Token</h3>
    <form onsubmit={(e) => { e.preventDefault(); creating = true; }} class="flex gap-3">
      <input type="text" bind:value={newName} placeholder="Token name (e.g., CI/CD Pipeline)" required
        class="flex-1 px-4 py-2 rounded-lg border text-sm focus:outline-none focus:ring-2"
        style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary);" />
      <button type="submit" disabled={creating}
        class="px-4 py-2 rounded-lg text-sm font-medium text-white disabled:opacity-50"
        style="background: var(--accent);">
        {creating ? 'Creating...' : 'Create Token'}
      </button>
    </form>
  </div>

  {#if newToken}
    <div class="p-4 rounded-xl border mb-6" style="background: rgba(34,197,94,0.1); border-color: rgba(34,197,94,0.3);">
      <p class="text-sm font-medium mb-2" style="color: var(--success);">✓ Token created! Copy it now — it won't be shown again.</p>
      <div class="flex items-center gap-2 p-2 rounded" style="background: var(--bg-tertiary);">
        <code class="flex-1 text-sm font-mono break-all">{newToken}</code>
        <button onclick={() => copyToken(newToken!)} class="px-3 py-1 rounded text-xs font-medium border" style="border-color: var(--border);">
          Copy
        </button>
      </div>
    </div>
  {/if}

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if tokens.length === 0}
    <div class="text-center py-16 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <svg class="w-12 h-12 mx-auto mb-4" style="color: var(--text-muted);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
      </svg>
      <p class="mb-2 font-medium">No service tokens</p>
      <p style="color: var(--text-secondary);">Create a token to access the API programmatically</p>
    </div>
  {:else}
    <div class="space-y-3">
      {#each tokens as token}
        <div class="flex items-center justify-between p-4 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <div>
            <p class="font-medium">{token.name}</p>
            <p class="text-sm font-mono" style="color: var(--text-muted);">{token.prefix || '***'}...</p>
            <p class="text-xs mt-1" style="color: var(--text-muted);">Created {new Date(token.created_at).toLocaleDateString()}</p>
          </div>
          <button class="p-2 rounded hover:bg-red-500/10" title="Revoke">
            <svg class="w-4 h-4" style="color: var(--danger);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>
