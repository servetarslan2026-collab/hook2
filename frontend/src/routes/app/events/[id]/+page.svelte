<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import type { Event } from '$lib/api/client';

  let event = $state<Event | null>(null);
  let loading = $state(true);
  let eventId = $derived($page.params.id);

  onMount(async () => {
    try {
      const { eventApi } = await import('$lib/api/client');
      event = await eventApi.get(eventId);
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  function formatJson(obj: object | string | undefined) {
    if (!obj) return '(empty)';
    if (typeof obj === 'string') {
      try { return JSON.stringify(JSON.parse(obj), null, 2); } catch { return obj; }
    }
    return JSON.stringify(obj, null, 2);
  }
</script>

<div class="p-8 max-w-4xl">
  <div class="mb-8">
    <div class="flex items-center gap-2 mb-2">
      <a href="/app/events" class="text-sm" style="color: var(--text-muted);">← Back to events</a>
    </div>
    <h1 class="text-2xl font-bold">Event Details</h1>
    <p class="mt-1 font-mono text-sm" style="color: var(--text-muted);">{eventId}</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if !event}
    <div class="text-center py-20 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <p class="text-lg font-medium">Event not found</p>
      <p class="mt-2" style="color: var(--text-secondary);">This event may have been cleaned up.</p>
    </div>
  {:else}
    <!-- Meta -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <p class="text-sm" style="color: var(--text-secondary);">Event Type</p>
        <div class="mt-2">
          <span class="inline-flex items-center px-2.5 py-1 rounded-lg text-sm font-mono font-medium" style="background: rgba(59,130,246,0.1); color: var(--accent);">
            {event.event_type}
          </span>
        </div>
      </div>
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <p class="text-sm" style="color: var(--text-secondary);">Application ID</p>
        <p class="text-sm font-mono mt-2">{event.application_id}</p>
      </div>
      <div class="p-5 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
        <p class="text-sm" style="color: var(--text-secondary);">Created</p>
        <p class="text-sm mt-2">{new Date(event.created_at).toLocaleString()}</p>
      </div>
    </div>

    <!-- Payload -->
    <div class="mb-6">
      <h3 class="text-sm font-medium mb-2" style="color: var(--text-secondary);">Payload</h3>
      <div class="rounded-xl border overflow-hidden" style="background: var(--bg-secondary); border-color: var(--border);">
        <pre class="p-4 text-sm font-mono overflow-x-auto" style="color: var(--text-primary);">{formatJson(event.payload)}</pre>
      </div>
    </div>

    <!-- Metadata -->
    {#if event.metadata}
      <div>
        <h3 class="text-sm font-medium mb-2" style="color: var(--text-secondary);">Metadata</h3>
        <div class="rounded-xl border overflow-hidden" style="background: var(--bg-secondary); border-color: var(--border);">
          <pre class="p-4 text-sm font-mono overflow-x-auto" style="color: var(--text-primary);">{formatJson(event.metadata)}</pre>
        </div>
      </div>
    {/if}
  {/if}
</div>
