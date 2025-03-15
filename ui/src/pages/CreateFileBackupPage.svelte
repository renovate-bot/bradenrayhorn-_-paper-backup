<script lang="ts">
  import QRCode from "qrcode";
  import PrintView from "../lib/PrintView.svelte";
  import { writeBarcode, type WriterOptions } from "zxing-wasm/writer";
  import { zxing } from "../zxing";

  let files = $state<FileList | null>(null);
  let passphrase = $state("");

  let error = $state("");

  let backup = $state<{ qr: string } | null>(null);

  async function generateQRCodes(text: string) {
    const chunks = [];
    for (let i = 0; i < text.length; i += 300) {
      chunks.push(text.slice(i, i + 300));
    }

    const qrCodes = [];

    for (const chunk of chunks) {
      const url = await QRCode.toDataURL(chunk, { errorCorrectionLevel: "M" });
      qrCodes.push(url);
    }

    return qrCodes;
  }

  async function doBackup() {
    if (files?.length !== 1 || passphrase.trim().length < 1) {
      error = "Please enter a single file and a passphrase.";
      return;
    }

    const bytes = await files[0].bytes();
    const res = paperBackup(bytes, passphrase);

    var result = zxing.generateBarcodeFromBinary(
      //res.qr,
      new Uint8Array([1, 2, 3, 4]),
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
