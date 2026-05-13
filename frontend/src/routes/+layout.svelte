<script lang="ts">
  import '../app.css';
  import { page } from '$app/stores';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  let { children } = $props();

  let sidebarOpen = $state(true);
  let currentOrg = $state<any>(null);
  let currentApp = $state<any>(null);

  const isAuthPage = $derived(
    ['/login', '/register', '/forgot-password'].includes($page.url.pathname)
  );

  function handleLogout() {
    auth.logout();
    goto('/login');
  }
</script>

{#if isAuthPage}
  <div class="min-h-screen flex items-center justify-center" style="background: var(--bg-primary);">
    {@render children()}
  </div>
{:else}
  <div class="flex h-screen overflow-hidden" style="background: var(--bg-primary);">
    <!-- Sidebar -->
    <aside class="flex flex-col w-64 border-r" style="background: var(--bg-secondary); border-color: var(--border);">
      <!-- Logo -->
      <div class="flex items-center gap-3 px-6 py-4 border-b" style="border-color: var(--border);">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center font-bold text-white" style="background: var(--accent);">
          W
        </div>
        <span class="text-lg font-semibold">Webhook</span>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
        <a href="/" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors
          {$page.url.pathname === '/' ? 'text-white' : ''}"
          style="{$page.url.pathname === '/' ? 'background: var(--bg-hover)' : ''} color: var(--text-secondary);">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
          </svg>
          Dashboard
        </a>

        <a href="/organizations" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
          style="{$page.url.pathname.startsWith('/organizations') ? 'background: var(--bg-hover); color: white;' : 'color: var(--text-secondary);'}">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
          </svg>
          Organizations
        </a>

        <a href="/settings" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
          style="{$page.url.pathname === '/settings' ? 'background: var(--bg-hover); color: white;' : 'color: var(--text-secondary);'}">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
          Settings
        </a>

        <div class="pt-4 mt-4 border-t" style="border-color: var(--border);">
          <p class="px-3 py-1 text-xs font-medium uppercase" style="color: var(--text-muted);">Application</p>
          <a href="/app/events" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm"
            style="color: var(--text-secondary);">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
            </svg>
            Events
          </a>
          <a href="/app/subscriptions" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm"
            style="color: var(--text-secondary);">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
            </svg>
            Subscriptions
          </a>
          <a href="/app/deliveries" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm"
            style="color: var(--text-secondary);">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"/>
            </svg>
            Deliveries
          </a>
          <a href="/app/secrets" class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm"
            style="color: var(--text-secondary);">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
            </svg>
            Secrets
          </a>
        </div>
      </nav>

      <!-- User section -->
      <div class="px-3 py-4 border-t" style="border-color: var(--border);">
        <div class="flex items-center gap-3 px-3 py-2">
          <div class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium text-white" style="background: var(--accent);">
            {$auth.user?.name?.[0] || 'U'}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate">{$auth.user?.name || 'User'}</p>
            <p class="text-xs truncate" style="color: var(--text-muted);">{$auth.user?.email || ''}</p>
          </div>
          <button onclick={handleLogout} class="p-1 rounded hover:bg-gray-800" title="Logout">
            <svg class="w-4 h-4" style="color: var(--text-muted);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
            </svg>
          </button>
        </div>
      </div>
    </aside>

    <!-- Main content -->
    <main class="flex-1 overflow-y-auto">
      {@render children()}
    </main>
  </div>
{/if}
