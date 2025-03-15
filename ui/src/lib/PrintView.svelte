<script lang="ts">
  const { qr }: { qr: string } = $props();
  console.log(qr);
</script>

<div class="qr">
  <img src={qr} alt="qr code" />
</div>

<style>
</style>
