<script lang="ts">
  import { onMount } from 'svelte';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  let name = $state('');
  let email = $state('');
  let currentPassword = $state('');
  let newPassword = $state('');
  let confirmPassword = $state('');
  let loading = $state(true);
  let saving = $state(false);
  let changingPassword = $state(false);
  let error = $state('');
  let success = $state('');
  let pwError = $state('');
  let pwSuccess = $state('');

  onMount(() => {
    if (!$auth.user) {
      goto('/login');
      return;
    }
    name = $auth.user.name || '';
    email = $auth.user.email || '';
    loading = false;
  });

  async function handleProfileUpdate() {
    error = '';
    success = '';
    saving = true;
    try {
      const { api } = await import('$lib/api/client');
      await api.put('users/me', { json: { name } }).json();
      success = 'Profile updated successfully';
    } catch (e: any) {
      error = e.message || 'Failed to update profile';
    } finally {
      saving = false;
    }
  }

  async function handlePasswordChange() {
    pwError = '';
    pwSuccess = '';

    if (newPassword !== confirmPassword) {
      pwError = 'Passwords do not match';
      return;
    }
    if (newPassword.length < 8) {
      pwError = 'Password must be at least 8 characters';
      return;
    }

    changingPassword = true;
    try {
      const { api } = await import('$lib/api/client');
      await api.put('users/me/password', { json: { current_password: currentPassword, new_password: newPassword } }).json();
      pwSuccess = 'Password changed successfully';
      currentPassword = '';
      newPassword = '';
      confirmPassword = '';
    } catch (e: any) {
      pwError = e.message || 'Failed to change password';
    } finally {
      changingPassword = false;
    }
  }
</script>

<div class="p-8 max-w-2xl">
  <div class="mb-8">
    <h1 class="text-2xl font-bold">Settings</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Manage your account and preferences</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else}
    <!-- Profile Section -->
    <div class="p-6 rounded-xl border mb-6" style="background: var(--bg-secondary); border-color: var(--border);">
      <h2 class="text-lg font-semibold mb-4">Profile</h2>
      <form onsubmit={(e) => { e.preventDefault(); handleProfileUpdate(); }} class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Name</label>
          <input type="text" bind:value={name} placeholder="Your name"
            class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
            style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Email</label>
          <input type="email" value={email} disabled
            class="w-full px-4 py-2.5 rounded-lg border text-sm opacity-60 cursor-not-allowed"
            style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary);" />
          <p class="mt-1 text-xs" style="color: var(--text-muted);">Email cannot be changed</p>
        </div>

        {#if error}
          <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(239,68,68,0.1); color: var(--danger);">{error}</div>
        {/if}
        {#if success}
          <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(34,197,94,0.1); color: var(--success);">{success}</div>
        {/if}

        <button type="submit" disabled={saving}
          class="px-4 py-2 rounded-lg text-sm font-medium text-white disabled:opacity-50"
          style="background: var(--accent);">
          {saving ? 'Saving...' : 'Save Profile'}
        </button>
      </form>
    </div>

    <!-- Password Section -->
    <div class="p-6 rounded-xl border mb-6" style="background: var(--bg-secondary); border-color: var(--border);">
      <h2 class="text-lg font-semibold mb-4">Change Password</h2>
      <form onsubmit={(e) => { e.preventDefault(); handlePasswordChange(); }} class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Current Password</label>
          <input type="password" bind:value={currentPassword} placeholder="••••••••" required
            class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
            style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">New Password</label>
          <input type="password" bind:value={newPassword} placeholder="••••••••" required minlength="8"
            class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
            style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Confirm New Password</label>
          <input type="password" bind:value={confirmPassword} placeholder="••••••••" required minlength="8"
            class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
            style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
        </div>

        {#if pwError}
          <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(239,68,68,0.1); color: var(--danger);">{pwError}</div>
        {/if}
        {#if pwSuccess}
          <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(34,197,94,0.1); color: var(--success);">{pwSuccess}</div>
        {/if}

        <button type="submit" disabled={changingPassword}
          class="px-4 py-2 rounded-lg text-sm font-medium text-white disabled:opacity-50"
          style="background: var(--accent);">
          {changingPassword ? 'Changing...' : 'Change Password'}
        </button>
      </form>
    </div>

    <!-- Account Info -->
    <div class="p-6 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <h2 class="text-lg font-semibold mb-4">Account</h2>
      <div class="space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-sm" style="color: var(--text-secondary);">Account created</span>
          <span class="text-sm">{new Date($auth.user?.created_at || '').toLocaleDateString()}</span>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm" style="color: var(--text-secondary);">User ID</span>
          <span class="text-sm font-mono" style="color: var(--text-muted);">{$auth.user?.id}</span>
        </div>
      </div>
    </div>
  {/if}
</div>
