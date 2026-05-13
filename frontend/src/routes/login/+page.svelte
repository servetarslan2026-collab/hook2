<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  let email = $state('');
  let password = $state('');
  let name = $state('');
  let isRegister = $state(false);
  let error = $state('');
  let loading = $state(false);

  async function handleSubmit() {
    error = '';
    loading = true;

    try {
      const { authApi } = await import('$lib/api/client');
      const result = isRegister
        ? await authApi.register({ email, password, name })
        : await authApi.login({ email, password });

      auth.login(result.user, result.token);
      goto('/');
    } catch (e: any) {
      error = e.message || 'Authentication failed';
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
    <h1 class="text-2xl font-bold">{isRegister ? 'Create account' : 'Welcome back'}</h1>
    <p class="mt-2" style="color: var(--text-secondary);">
      {isRegister ? 'Get started with webhooks' : 'Sign in to your account'}
    </p>
  </div>

  <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
    {#if isRegister}
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Name</label>
        <input type="text" bind:value={name} placeholder="Your name"
          class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
      </div>
    {/if}

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
    </div>

    {#if error}
      <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(239,68,68,0.1); color: var(--danger); border: 1px solid rgba(239,68,68,0.2);">
        {error}
      </div>
    {/if}

    <button type="submit" disabled={loading}
      class="w-full py-2.5 rounded-lg text-sm font-medium text-white transition-colors disabled:opacity-50"
      style="background: var(--accent);">
      {loading ? 'Loading...' : (isRegister ? 'Create account' : 'Sign in')}
    </button>
  </form>

  <p class="mt-6 text-center text-sm" style="color: var(--text-secondary);">
    {isRegister ? 'Already have an account?' : "Don't have an account?"}
    <button onclick={() => { isRegister = !isRegister; error = ''; }} class="font-medium" style="color: var(--accent);">
      {isRegister ? 'Sign in' : 'Sign up'}
    </button>
  </p>
</div>
