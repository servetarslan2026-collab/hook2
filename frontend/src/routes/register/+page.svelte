<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  let name = $state('');
  let email = $state('');
  let password = $state('');
  let error = $state('');
  let loading = $state(false);

  async function handleSubmit() {
    error = '';
    loading = true;
    try {
      const { authApi } = await import('$lib/api/client');
      const result = await authApi.register({ email, password, name });
      auth.login(result.user, result.token);
      goto('/');
    } catch (e: any) {
      error = e.message || 'Registration failed';
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
    <h1 class="text-2xl font-bold">Create account</h1>
    <p class="mt-2" style="color: var(--text-secondary);">Get started with webhooks</p>
  </div>

  <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
    <div>
      <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Name</label>
      <input type="text" bind:value={name} placeholder="Your name" required
        class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
        style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
    </div>

    <div>
      <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Email</label>
      <input type="email" bind:value={email} placeholder="you@example.com" required
        class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
        style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
    </div>

    <div>
      <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Password</label>
      <input type="password" bind:value={password} placeholder="••••••••" required minlength="8"
        class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
        style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
      <p class="mt-1 text-xs" style="color: var(--text-muted);">At least 8 characters</p>
    </div>

    {#if error}
      <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(239,68,68,0.1); color: var(--danger); border: 1px solid rgba(239,68,68,0.2);">
        {error}
      </div>
    {/if}

    <button type="submit" disabled={loading}
      class="w-full py-2.5 rounded-lg text-sm font-medium text-white transition-colors disabled:opacity-50"
      style="background: var(--accent);">
      {loading ? 'Creating account...' : 'Create account'}
    </button>
  </form>

  <p class="mt-6 text-center text-sm" style="color: var(--text-secondary);">
    Already have an account?
    <a href="/login" class="font-medium" style="color: var(--accent);">Sign in</a>
  </p>
</div>
