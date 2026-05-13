<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';

  let eventType = $state('');
  let payloadText = $state('{\n  "message": "Hello World"\n}');
  let metadataText = $state('');
  let loading = $state(false);
  let error = $state('');
  let success = $state(false);
  let appId = $derived($page.params.id);

  async function sendEvent() {
    error = '';
    success = false;
    loading = true;

    try {
      const payload = JSON.parse(payloadText);
      const metadata = metadataText ? JSON.parse(metadataText) : undefined;

      const { eventApi } = await import('$lib/api/client');
      await eventApi.send(appId, { event_type: eventType, payload, metadata });
      success = true;
      setTimeout(() => { success = false; }, 3000);
    } catch (e: any) {
      error = e.message || 'Failed to send event';
    } finally {
      loading = false;
    }
  }
</script>

<div class="p-8 max-w-2xl">
  <div class="mb-8">
    <h1 class="text-2xl font-bold">Send Event</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Send a test event to trigger webhook deliveries</p>
  </div>

  <form onsubmit={(e) => { e.preventDefault(); sendEvent(); }} class="space-y-6">
    <div>
      <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Event Type</label>
      <input type="text" bind:value={eventType} placeholder="e.g., order.created, user.signed_up" required
        class="w-full px-4 py-2.5 rounded-lg border text-sm focus:outline-none focus:ring-2"
        style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
    </div>

    <div>
      <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Payload (JSON)</label>
      <textarea bind:value={payloadText} rows="10" required
        class="w-full px-4 py-2.5 rounded-lg border text-sm font-mono focus:outline-none focus:ring-2"
        style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);"></textarea>
    </div>

    <div>
      <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Metadata (JSON, optional)</label>
      <textarea bind:value={metadataText} rows="4"
        placeholder='{"source": "test", "user_id": "123"}'
        class="w-full px-4 py-2.5 rounded-lg border text-sm font-mono focus:outline-none focus:ring-2"
        style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);"></textarea>
    </div>

    {#if error}
      <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(239,68,68,0.1); color: var(--danger);">
        {error}
      </div>
    {/if}

    {#if success}
      <div class="px-4 py-3 rounded-lg text-sm" style="background: rgba(34,197,94,0.1); color: var(--success);">
        ✓ Event sent successfully!
      </div>
    {/if}

    <div class="flex gap-3">
      <button type="submit" disabled={loading}
        class="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium text-white disabled:opacity-50"
        style="background: var(--accent);">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
        </svg>
        {loading ? 'Sending...' : 'Send Event'}
      </button>
      <a href="/app/events" class="px-4 py-2 rounded-lg text-sm font-medium border" style="color: var(--text-secondary); border-color: var(--border);">
        Back to Events
      </a>
    </div>
  </form>

  <!-- API Example -->
  <div class="mt-8 p-4 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
    <h3 class="text-sm font-medium mb-3">Or send via API</h3>
    <pre class="p-3 rounded text-sm font-mono overflow-x-auto" style="background: var(--bg-tertiary); color: var(--text-secondary);">{`curl -X POST https://your-domain.com/api/v1/applications/{app_id}/events \\
  -H "Content-Type: application/json" \\
  -H "X-API-Key: your_api_key" \\
  -d '{
    "event_type": "order.created",
    "payload": {
      "order_id": "12345",
      "amount": 99.99
    }
  }'`}</pre>
  </div>
</div>
