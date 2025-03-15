<script lang="ts">
  import NotFoundPageState from "./lib/NotFoundPageState.svelte";
  import Router from "./lib/router/Router.svelte";
  import type { Route } from "./lib/router/Router.svelte";
  import CreateFileBackupPage from "./pages/CreateFileBackupPage.svelte";
  import DecryptFileBackupPage from "./pages/DecryptFileBackupPage.svelte";

  const routes: Route[] = [
    { path: "file/create", component: CreateFileBackupPage },
    { path: "file/decrypt", component: DecryptFileBackupPage },
    { path: ".*", component: NotFoundPageState },
  ];
</script>

<main>
  <Router {routes} />
</main>
