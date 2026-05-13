<script lang="ts">
  let payload = $state('{\n  "event_type": "order.created",\n  "payload": {\n    "order_id": "12345",\n    "amount": 99.99\n  }\n}');
  let secret = $state('');
  let signature = $state('');
  let computed = $state(false);
  let testSignature = $state('');
  let verifyResult = $state<'match' | 'mismatch' | null>(null);

  async function computeSignature() {
    if (!secret || !payload) return;

    const encoder = new TextEncoder();
    const keyData = encoder.encode(secret);
    const data = encoder.encode(payload);

    const cryptoKey = await crypto.subtle.importKey(
      'raw', keyData, { name: 'HMAC', hash: 'SHA-256' }, false, ['sign']
    );

    const sig = await crypto.subtle.sign('HMAC', cryptoKey, data);
    const hex = Array.from(new Uint8Array(sig)).map(b => b.toString(16).padStart(2, '0')).join('');
    signature = `sha256=${hex}`;
    computed = true;
    verifyResult = null;
  }

  function verifySignature() {
    if (!signature || !testSignature) return;
    verifyResult = signature === testSignature.trim() ? 'match' : 'mismatch';
  }

  function copySignature() {
    navigator.clipboard.writeText(signature);
  }
</script>

<div class="p-8 max-w-3xl">
  <div class="mb-8">
    <h1 class="text-2xl font-bold">Webhook Signature Test</h1>
    <p class="mt-1" style="color: var(--text-secondary);">Compute and verify HMAC-SHA256 webhook signatures</p>
  </div>

  <!-- Compute Section -->
  <div class="p-6 rounded-xl border mb-6" style="background: var(--bg-secondary); border-color: var(--border);">
    <h2 class="text-lg font-semibold mb-4">Compute Signature</h2>
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Payload (JSON body)</label>
        <textarea bind:value={payload} rows="8" required
          class="w-full px-4 py-2.5 rounded-lg border text-sm font-mono focus:outline-none focus:ring-2 resize-none"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);"></textarea>
      </div>
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Secret (webhook signing secret)</label>
        <input type="text" bind:value={secret} placeholder="whsec_sub_..." required
          class="w-full px-4 py-2.5 rounded-lg border text-sm font-mono focus:outline-none focus:ring-2"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
      </div>

      <button onclick={computeSignature}
        class="px-4 py-2 rounded-lg text-sm font-medium text-white"
        style="background: var(--accent);">
        Compute Signature
      </button>

      {#if computed}
        <div class="p-4 rounded-lg border" style="background: var(--bg-tertiary); border-color: var(--border);">
          <div class="flex items-center justify-between mb-2">
            <p class="text-sm font-medium" style="color: var(--text-secondary);">X-Webhook-Signature</p>
            <button onclick={copySignature} class="px-2 py-0.5 rounded text-xs border" style="border-color: var(--border); color: var(--text-secondary);">
              Copy
            </button>
          </div>
          <code class="text-sm font-mono break-all" style="color: var(--success);">{signature}</code>
        </div>
      {/if}
    </div>
  </div>

  <!-- Verify Section -->
  <div class="p-6 rounded-xl border mb-6" style="background: var(--bg-secondary); border-color: var(--border);">
    <h2 class="text-lg font-semibold mb-4">Verify Signature</h2>
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium mb-1.5" style="color: var(--text-secondary);">Signature to verify</label>
        <input type="text" bind:value={testSignature} placeholder="sha256=..."
          class="w-full px-4 py-2.5 rounded-lg border text-sm font-mono focus:outline-none focus:ring-2"
          style="background: var(--bg-tertiary); border-color: var(--border); color: var(--text-primary); --tw-ring-color: var(--accent);" />
      </div>

      <button onclick={verifySignature}
        class="px-4 py-2 rounded-lg text-sm font-medium text-white"
        style="background: var(--accent);">
        Verify
      </button>

      {#if verifyResult}
        <div class="p-3 rounded-lg text-sm font-medium"
          style="background: {verifyResult === 'match' ? 'rgba(34,197,94,0.1)' : 'rgba(239,68,68,0.1)'}; color: {verifyResult === 'match' ? 'var(--success)' : 'var(--danger)'};">
          {#if verifyResult === 'match'}
            ✓ Signatures match! The payload is authentic.
          {:else}
            ✗ Signatures do not match. The payload may have been tampered with.
          {/if}
        </div>
      {/if}
    </div>
  </div>

  <!-- How it works -->
  <div class="p-6 rounded-xl border" style="background: var(--bg-secondary); border-color: var(--border);">
    <h2 class="text-lg font-semibold mb-3">How it works</h2>
    <div class="space-y-3 text-sm" style="color: var(--text-secondary);">
      <p>Every webhook request includes an <code class="font-mono" style="color: var(--accent);">X-Webhook-Signature</code> header containing an HMAC-SHA256 signature of the request body.</p>
      <p>To verify:</p>
      <ol class="list-decimal list-inside space-y-1 ml-2">
        <li>Compute HMAC-SHA256 of the raw request body using your webhook secret</li>
        <li>Compare the result with the <code class="font-mono" style="color: var(--accent);">sha256=...</code> value from the header</li>
        <li>If they match, the payload is authentic and hasn't been tampered with</li>
      </ol>
      <div class="mt-4 p-3 rounded-lg font-mono text-xs" style="background: var(--bg-tertiary);">
        <p style="color: var(--text-muted);"># Python verification example</p>
        <p>import hmac, hashlib</p>
        <p>expected = hmac.new(secret.encode(), payload, hashlib.sha256).hexdigest()</p>
        <p>hmac.compare_digest(f"sha256={{expected}}", signature)</p>
      </div>
    </div>
  </div>
</div>
