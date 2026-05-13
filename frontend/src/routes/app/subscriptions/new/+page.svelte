<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import type { Application, EventType } from '$lib/api/client';

  let apps = $state<Application[]>([]);
  let eventTypes = $state<EventType[]>([]);
  let selectedAppId = $state('');
  let targetUrl = $state('');
  let description = $state('');
  let selectedEventTypes = $state<string[]>([]);
  let loading = $state(true);
  let saving = $state(false);
  let error = $state('');

  onMount(async () => {
    try {
      const { orgApi, appApi } = await import('$lib/api/client');
      const orgs = await orgApi.list();
      if (orgs.length > 0) {
        apps = await appApi.list(orgs[0].id);
        if (apps.length > 0) {
          selectedAppId = apps[0].id;
          await loadEventTypes(selectedAppId);
        }
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  });

  async function loadEventTypes(appId: string) {
    try {
      const { eventTypeApi } = await import('$lib/api/client');
      eventTypes = await eventTypeApi.list(appId);
    } catch (e) {
      console.error(e);
    }
  }

  async function handleAppChange() {
    selectedEventTypes = [];
    if (selectedAppId) {
      await loadEventTypes(selectedAppId);
    }
  }

  function toggleEventType(name: string) {
    if (selectedEventTypes.includes(name)) {
      selectedEventTypes = selectedEventTypes.filter(t => t !== name);
    } else {
      selectedEventTypes = [...selectedEventTypes, name];
    }
  }

  async function handleSubmit() {
    if (!selectedAppId || !targetUrl || selectedEventTypes.length === 0) {
      error = 'Please fill in all required fields and select at least one event type';
      return;
    }
    error = '';
    saving = true;
    try {
      const { subscriptionApi } = await import('$lib/api/client');
      await subscriptionApi.create(selectedAppId, {
        event_types: selectedEventTypes,
        target_url: targetUrl,
        description
      });
      goto('/app/subscriptions');
    } catch (e: any) {
      error = e.message || 'Failed to create subscription';
    } finally {
      saving = false;
    }
  }
</script>

<div class="p-8 max-w-2xl">
  <div class="mb-8">
    <div class="flex items-center gap-2 mb-2">
      <a href="/app/subscriptions" class="text-sm" style="color: var(--text-muted);">← Back to subscriptions</a>
    </div>
    <h1 class="text-2xl font-bold">New Subscription</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Create a new webhook endpoint to receive events</p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 rounded-full animate-spin" style="border-color: var(--border); border-top-color: var(--accent);"></div>
    </div>
  {:else}
    <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-6">
      <!-- App Selection -->
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Application *</label>
        <select bind:value={selectedAppId} onchange={handleAppChange}
          class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);">
          {#each apps as app}
            <option value={app.id}>{app.name}</option>
          {/each}
        </select>
      </div>

      <!-- Target URL -->
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Webhook URL *</label>
        <input type="url" bind:value={targetUrl} placeholder="https://your-app.com/webhooks" required
          class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
        <p class="mt-1 text-xs" style="color: var(--text-muted);">We'll send POST requests to this URL when events match</p>
      </div>

      <!-- Description -->
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Description</label>
        <input type="text" bind:value={description} placeholder="e.g., Production order notifications"
          class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
      </div>

      <!-- Event Types -->
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Event Types *</label>
        {#if eventTypes.length === 0}
          <div class="p-4 rounded-xl border text-center" style="background: var(--bg-secondary); border-color: var(--border);">
            <p class="text-sm" style="color: var(--text-muted);">No event types defined. Create event types first.</p>
            <a href="/app/{selectedAppId}/event-types" class="text-sm font-medium mt-1 inline-block" style="color: var(--accent);">Create event types →</a>
          </div>
        {:else}
          <div class="space-y-2">
            {#each eventTypes as et}
              <label class="flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-colors"
                style="background: {selectedEventTypes.includes(et.name) ? 'rgba(59,130,246,0.1)' : 'var(--bg-secondary)'}; border-color: {selectedEventTypes.includes(et.name) ? 'var(--accent)' : 'var(--border)'};">
                <input type="checkbox" checked={selectedEventTypes.includes(et.name)} onchange={() => toggleEventType(et.name)}
                  class="w-4 h-4 rounded" style="accent-color: var(--accent);" />
                <div>
                  <span class="text-sm font-mono font-medium">{et.name}</span>
                  {#if et.description}
                    <span class="text-sm ml-2" style="color: var(--text-muted);">— {et.description}</span>
                  {/if}
                </div>
              </label>
            {/each}
          </div>
        {/if}
      </div>

      {#if error}
        <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(239,68,68,0.1); color: var(--danger);">
          {error}
        </div>
      {/if}

      <!-- Actions -->
      <div class="flex gap-3">
        <button type="submit" disabled={saving}
          class="px-6 py-2.5 rounded-lg text-sm font-medium text-white disabled:opacity-50"
          style="background: var(--accent);">
          {saving ? 'Creating...' : 'Create Subscription'}
        </button>
        <a href="/app/subscriptions" class="px-4 py-2.5 rounded-lg text-sm font-medium border" style="color: var(--text-secondary); border-color: var(--border);">
          Cancel
        </a>
      </div>
    </form>
  {/if}
</div>
