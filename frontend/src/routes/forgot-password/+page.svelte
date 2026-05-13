<script lang="ts">
  let email = $state('');
  let loading = $state(false);
  let sent = $state(false);
  let error = $state('');

  async function handleSubmit() {
    error = '';
    loading = true;
    try {
      const { authApi } = await import('$lib/api/client');
      // Call forgot-password endpoint
      await fetch(`${import.meta.env.PUBLIC_API_URL || '/api/v1'}/auth/forgot-password`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
      });
      sent = true;
    } catch (e: any) {
      error = e.message || 'Failed to send reset email';
    } finally {
      loading = false;
    }
  }
</script>

<div class="w-full max-w-md p-8">
  <div class="text-center mb-8">
    <div class="w-12 h-12 rounded-xl flex items-center justify-center font-bold text-white text-xl mx-auto mb-4" style="background: var(--accent);">
      W
    </div>
    <h1 class="text-2xl font-bold">Reset password</h1>
    <p class="mt-2" style="color: var(--text-secondary);">Enter your email to receive a reset link</p>
  </div>

  {#if sent}
    <div class="text-center py-8">
      <div class="w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4" style="background: rgba(34,197,94,0.1);">
        <svg class="w-8 h-8" style="color: var(--success);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
        </svg>
      </div>
      <h2 class="text-lg font-semibold mb-2">Check your email</h2>
      <p class="mb-6" style="color: var(--text-secondary);">
        If an account exists for <strong>{email}</strong>, we've sent a password reset link.
      </p>
      <a href="/login" class="text-sm font-medium" style="color: var(--accent);">← Back to sign in</a>
    </div>
  {:else}
    <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Email</label>
        <input type="email" bind:value={email} placeholder="you@example.com" required
          class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
      </div>

      {#if error}
        <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(239,68,68,0.1); color: var(--danger); border: 1px solid rgba(239,68,68,0.2);">
          {error}
        </div>
      {/if}

      <button type="submit" disabled={loading}
        class="w-full py-2.5 rounded-lg text-sm font-medium text-white transition-colors disabled:opacity-50"
        style="background: var(--accent);">
        {loading ? 'Sending...' : 'Send reset link'}
      </button>
    </form>

    <p class="mt-6 text-center text-sm" style="color: var(--text-secondary);">
      Remember your password?
      <a href="/login" class="font-medium" style="color: var(--accent);">Sign in</a>
    </p>
  {/if}
</div>
