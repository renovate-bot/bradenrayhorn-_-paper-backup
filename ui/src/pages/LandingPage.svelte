<script lang="ts">
  import { link } from "../lib/router/use-link.svelte";
</script>

Pick a tool.

<ul>
  <li><a use:link href="/#/backup/file">Create File Backup</a></li>
  <li><a use:link href="/#/restore/file">Restore File Backup</a></li>
</ul>
