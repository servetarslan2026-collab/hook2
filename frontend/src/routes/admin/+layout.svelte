<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { auth } from '$lib/stores/auth';

  let { children } = $props();
  let isAdmin = $state(false);
  let loading = $state(true);

  onMount(() => {
    const unsub = auth.subscribe(state => {
      if (!state.token) {
        goto('/login');
        return;
      }
      // Check admin status from user object
      isAdmin = (state.user as any)?.is_admin === true;
      if (!isAdmin && !loading) {
        goto('/');
      }
      loading = false;
    });
    return unsub;
  });
</script>

{#if loading}
  <div class="flex items-center justify-center h-screen" style="background: var(--bg-primary);">
    <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
  </div>
{:else if isAdmin}
  <div class="flex h-screen overflow-hidden" style="background: var(--bg-primary);">
    <!-- Admin Sidebar -->
    <aside class="flex flex-col w-64 border-r" style="background: var(--bg-secondary); border-color: var(--border);">
      <div class="flex items-center gap-3 px-6 py-4 border-b" style="border-color: var(--border);">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center font-bold text-white" style="background: var(--danger);">
          A
        </div>
        <div>
          <span class="text-lg font-semibold">Admin</span>
          <p class="text-xs" style="color: var(--text-muted);">System Management</p>
        </div>
      </div>

      <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
        <a href="/admin" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
          style="{$page.url.pathname === '/admin' ? 'background: var(--bg-hover); color: white;' : 'color: var(--text-secondary);'}">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
          </svg>
          Dashboard
        </a>
        <a href="/admin/users" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
          style="{$page.url.pathname.startsWith('/admin/users') ? 'background: var(--bg-hover); color: white;' : 'color: var(--text-secondary);'}">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"/>
          </svg>
          Users
        </a>
        <a href="/admin/organizations" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
          style="{$page.url.pathname.startsWith('/admin/organizations') ? 'background: var(--bg-hover); color: white;' : 'color: var(--text-secondary);'}">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
          </svg>
          Organizations
        </a>
        <a href="/admin/dead-letters" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
          style="{$page.url.pathname.startsWith('/admin/dead-letters') ? 'background: var(--bg-hover); color: white;' : 'color: var(--text-secondary);'}">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z"/>
          </svg>
          Dead Letters
        </a>
      </nav>

      <div class="px-3 py-4 border-t" style="border-color: var(--border);">
        <a href="/" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm"
          style="color: var(--text-secondary);">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 17l-5-5m0 0l5-5m-5 5h12"/>
          </svg>
          Back to App
        </a>
      </div>
    </aside>

    <main class="flex-1 overflow-y-auto">
      {@render children()}
    </main>
  </div>
{:else}
  <div class="flex items-center justify-center h-screen" style="background: var(--bg-primary);">
    <div class="text-center">
      <h1 class="text-2xl font-bold mb-2">Access Denied</h1>
      <p style="color: var(--text-secondary);">You don't have admin privileges.</p>
      <a href="/" class="inline-block mt-4 px-4 py-2 rounded-lg text-sm text-white" style="background: var(--accent);">Go Home</a>
    </div>
  </div>
{/if}
