<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';

  let loading = $state(true);
  let error = $state('');
  let success = $state(false);
  let orgName = $state('');
  let token = $derived($page.params.token);

  onMount(async () => {
    if (!$auth.user) {
      goto(`/login?redirect=/invite/${token}`);
      return;
    }

    try {
      const { orgApi } = await import('$lib/api/client');
      await orgApi.acceptInvitation(token);
      success = true;
      setTimeout(() => goto('/organizations'), 2000);
    } catch (e: any) {
      error = e.message || 'Failed to accept invitation';
    } finally {
      loading = false;
    }
  });
</script>

<div class="w-full max-w-md p-8 text-center">
  {#if loading}
    <div class="w-12 h-12 border-2 rounded-full animate-spin mx-auto mb-4" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    <h1 class="text-xl font-semibold">Accepting invitation...</h1>
  {:else if success}
    <div class="w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4" style="background: rgba(34,197,94,0.1);">
      <svg class="w-8 h-8" style="color: var(--success);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
      </svg>
    </div>
    <h1 class="text-2xl font-bold mb-2">Invitation accepted!</h1>
    <p style="color: var(--text-secondary);">Redirecting to organizations...</p>
  {:else}
    <div class="w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4" style="background: rgba(239,68,68,0.1);">
      <svg class="w-8 h-8" style="color: var(--danger);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
      </svg>
    </div>
    <h1 class="text-2xl font-bold mb-2">Invitation failed</h1>
    <p class="mb-6" style="color: var(--text-secondary);">{error}</p>
    <a href="/" class="text-sm font-medium" style="color: var(--accent);">Go to dashboard →</a>
  {/if}
</div>
