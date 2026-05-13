<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';

  let loading = $state(false);
  let sent = $state(false);
  let error = $state('');
  let devToken = $state('');

  onMount(() => {
    if (!$auth.user) goto('/login');
  });

  async function sendVerification() {
    error = '';
    loading = true;
    try {
      const { authApi } = await import('$lib/api/client');
      const result = await authApi.sendVerification();
      sent = true;
      if (result.token) devToken = result.token;
    } catch (e: any) {
      error = e.message || 'Failed to send verification email';
    } finally {
      loading = false;
    }
  }
</script>

<div class="w-full max-w-md p-8 text-center">
  <div class="w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-6" style="background: rgba(59,130,246,0.1);">
    <svg class="w-8 h-8" style="color: var(--accent);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
    </svg>
  </div>

  {#if sent}
    <h1 class="text-2xl font-bold mb-2">Check your email</h1>
    <p class="mb-6" style="color: var(--text-secondary);">
      We've sent a verification link to <strong>{$auth.user?.email}</strong>
    </p>
    {#if devToken}
      <div class="p-3 rounded-lg border text-left mb-4" style="background: var(--bg-secondary); border-color: var(--border);">
        <p class="text-xs mb-1" style="color: var(--warning);">⚠ Dev mode — token visible:</p>
        <a href="/verify/{devToken}" class="text-sm font-mono break-all" style="color: var(--accent);">{devToken}</a>
      </div>
    {/if}
    <button onclick={() => { sent = false; sendVerification(); }}
      class="text-sm" style="color: var(--accent);">
      Didn't receive it? Send again
    </button>
  {:else}
    <h1 class="text-2xl font-bold mb-2">Verify your email</h1>
    <p class="mb-6" style="color: var(--text-secondary);">
      Verify your email address to secure your account
    </p>

    {#if error}
      <div class="px-4 py-3 rounded-lg text-sm mb-4" style="background: rgba(239,68,68,0.1); color: var(--danger);">{error}</div>
    {/if}

    <button onclick={sendVerification} disabled={loading}
      class="px-6 py-2.5 rounded-lg text-sm font-medium text-white disabled:opacity-50"
      style="background: var(--accent);">
      {loading ? 'Sending...' : 'Send verification email'}
    </button>
  {/if}

  <p class="mt-6">
    <a href="/" class="text-sm" style="color: var(--text-muted);">← Back to dashboard</a>
  </p>
</div>
