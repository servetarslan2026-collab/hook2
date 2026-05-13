<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import type { EventType } from '$lib/api/client';

  let eventTypes = $state<EventType[]>([]);
  let loading = $state(true);
  let newName = $state('');
  let newDescription = $state('');
  let creating = $state(false);
  let appId = $derived($page.params.id);

  onMount(async () => {
    await loadEventTypes();
  });

  async function loadEventTypes() {
    try {
      const { eventTypeApi } = await import('$lib/api/client');
      eventTypes = await eventTypeApi.list(appId);
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function createEventType() {
    if (!newName.trim()) return;
    creating = true;
    try {
      const { eventTypeApi } = await import('$lib/api/client');
      const et = await eventTypeApi.create(appId, { name: newName, description: newDescription });
      eventTypes = [...eventTypes, et];
      newName = '';
      newDescription = '';
    } catch (e: any) {
      console.error(e);
    } finally {
      creating = false;
    }
  }

  async function deleteEventType(id: string) {
    if (!confirm('Delete this event type?')) return;
    try {
      const { eventTypeApi } = await import('$lib/api/client');
      await eventTypeApi.delete(appId, id);
      eventTypes = eventTypes.filter(et => et.id !== id);
    } catch (e) {
      console.error(e);
    }
  }
</script>

<div class="p-8">
  <div class="flex items-center justify-between mb-8">
    <div>
      <h1 class="text-2xl font-bold">Event Types</h1>
      <p class="mt-1" style="color: var(--text-secondary);">Define the types of events your application can emit</p>
    </div>
  </div>

  <!-- Create Form -->
  <div class="p-4 rounded-xl border mb-6" style="background: var(--bg-secondary); border-color: var(--border);">
    <h3 class="text-sm font-medium mb-3">Create Event Type</h3>
    <form onsubmit={(e) => { e.preventDefault(); createEventType(); }} class="space-y-3">
      <div class="flex gap-3">
        <input type="text" bind:value={newName} placeholder="e.g., order.created" required
          class="flex-1 px-4 py-2 rounded-lg border text-sm focus:outline-none focus:ring-2"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary);" />
        <input type="text" bind:value={newDescription} placeholder="Description (optional)"
          class="flex-1 px-4 py-2 rounded-lg border text-sm focus:outline-none focus:ring-2"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary);" />
        <button type="submit" disabled={creating}
          class="px-4 py-2 rounded-lg text-sm font-medium text-white disabled:opacity-50"
          style="background: var(--accent);">
          {creating ? 'Creating...' : 'Create'}
        </button>
      </div>
    </form>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else if eventTypes.length === 0}
    <div class="text-center py-16 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
      <svg class="w-12 h-12 mx-auto mb-4" style="color: var(--text-muted);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"/>
      </svg>
      <p class="mb-2 font-medium">No event types</p>
      <p style="color: var(--text-secondary);">Define event types like order.created, user.signed_up, etc.</p>
    </div>
  {:else}
    <div class="space-y-3">
      {#each eventTypes as et}
        <div class="flex items-center justify-between p-4 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-mono font-medium" style="background: rgba(59,130,246,0.1); color: var(--accent);">
              {et.name}
            </span>
            <span class="text-sm" style="color: var(--text-secondary);">{et.description || 'No description'}</span>
          </div>
          <button onclick={() => deleteEventType(et.id)} class="p-1.5 rounded hover:bg-red-500/10" title="Delete">
            <svg class="w-4 h-4" style="color: var(--danger);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>
