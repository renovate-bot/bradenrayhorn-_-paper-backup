<script lang="ts">
  import { onDestroy } from "svelte";
  import PrintView from "../lib/PrintView.svelte";
  import { zxing } from "../zxing";

  let files = $state<FileList | null>(null);
  let passphrase = $state("");

  let error = $state("");

  let backup = $state<{ qr: string } | null>(null);

  async function doBackup() {
    if (files?.length !== 1 || passphrase.trim().length < 1) {
      error = "Please enter a single file and a passphrase.";
      return;
    }

    const bytes = await files[0].bytes();
    const res = window.paperBackup(bytes, passphrase);

    var result = zxing.generateBarcodeFromBinary(
      res.qr,
      "QRCODE",
      "BINARY",
      0,
      500,
      500,
      5,
    );

    const file = URL.createObjectURL(new Blob([result.image]));
    backup = { qr: file };

    result.delete();
  }

  // cleanup QR code data
  onDestroy(() => {
    if (backup?.qr) {
      URL.revokeObjectURL(backup.qr);
    }
  });
</script>

{#if backup !== null}
  <PrintView qr={backup.qr} />
{:else}
  <div>
    Create a file backup.

    <div>
      <label
        >Enter a passphrase
        <input bind:value={passphrase} name="passphrase" type="text" /></label
      >
    </div>

    <div>
      <label
        >Upload a file:
        <input bind:files name="file" type="file" /></label
      >
    </div>

    <button onclick={doBackup}>Backup</button>

    <span>{error}</span>
  </div>
{/if}
