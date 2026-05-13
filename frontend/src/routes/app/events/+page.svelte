<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import type { Event, PaginatedResponse } from '$lib/api/client';

  let events = $state<Event[]>([]);
  let total = $state(0);
  let currentPage = $state(1);
  let loading = $state(true);
  let appId = $derived($page.params.id);

  onMount(async () => {
    await loadEvents();
  });

  async function loadEvents(page = 1) {
    loading = true;
    try {
      const { eventApi } = await import('$lib/api/client');
      const result = await eventApi.list(appId, { page, per_page: 20 });
      events = result.data;
      total = result.total;
      currentPage = page;
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }
</script>

<div class="p-8">
  <div class="flex items-center justify-between mb-8">
    <div>
      <h1 class="text-2xl font-bold">Events</h1>
      <p class="mt-1" style="color: var(--text-secondary);">Inbound events for this application</p>
    </div>
    <a href="/app/events/send"
      class="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium text-white"
      style="background: var(--accent);">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
      </svg>
      Send Event
    </a>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if events.length === 0}
    <div class="text-center py-20 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <svg class="w-12 h-12 mx-auto mb-4" style="color: var(--text-muted);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
      </svg>
      <p class="mb-2 font-medium">No events yet</p>
      <p class="mb-4" style="color: var(--text-secondary);">Send your first event via the API or test button</p>
      <a href="/app/events/send" class="text-sm font-medium" style="color: var(--accent);">Send test event →</a>
    </div>
  {:else}
    <div class="rounded-xl border overflow-hidden" style="background: var(--bg-secondary); border-color: var(--border);">
      <table class="w-full">
        <thead>
          <tr style="background: var(--bg-tertiary);">
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Event Type</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">ID</th>
            <th class="px-4 py-3 text-left text-xs font-medium uppercase" style="color: var(--text-muted);">Created</th>
            <th class="px-4 py-3 text-right text-xs font-medium uppercase" style="color: var(--text-muted);">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y" style="border-color: var(--border);">
          {#each events as event}
            <tr class="hover:bg-white/5 transition-colors">
              <td class="px-4 py-3">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium" style="background: rgba(59,130,246,0.1); color: var(--accent);">
                  {event.event_type}
                </span>
              </td>
              <td class="px-4 py-3 text-sm font-mono" style="color: var(--text-muted);">
                {event.id.slice(0, 8)}...
              </td>
              <td class="px-4 py-3 text-sm" style="color: var(--text-secondary);">
                {new Date(event.created_at).toLocaleString()}
              </td>
              <td class="px-4 py-3 text-right">
                <a href="/app/events/{event.id}" class="text-sm font-medium" style="color: var(--accent);">View →</a>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    {#if total > 20}
      <div class="flex items-center justify-between mt-4">
        <p class="text-sm" style="color: var(--text-muted);">Showing {(currentPage - 1) * 20 + 1}-{Math.min(currentPage * 20, total)} of {total}</p>
        <div class="flex gap-2">
          <button onclick={() => loadEvents(currentPage - 1)} disabled={currentPage === 1}
            class="px-3 py-1 rounded text-sm border disabled:opacity-50" style="border-color: var(--border);">
            Previous
          </button>
          <button onclick={() => loadEvents(currentPage + 1)} disabled={currentPage * 20 >= total}
            class="px-3 py-1 rounded text-sm border disabled:opacity-50" style="border-color: var(--border);">
            Next
          </button>
        </div>
      </div>
    {/if}
  {/if}
</div>
