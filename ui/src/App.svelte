<script lang="ts">
  import NotFoundPageState from "./lib/NotFoundPageState.svelte";
  import Router from "./lib/router/Router.svelte";
  import type { Route } from "./lib/router/Router.svelte";
  import CreateFileBackupPage from "./pages/CreateFileBackupPage.svelte";
  import DecryptFileBackupPage from "./pages/DecryptFileBackupPage.svelte";
  import LandingPage from "./pages/LandingPage.svelte";

  const routes: Route[] = [
    { path: "/", component: LandingPage },
    { path: "backup/file", component: CreateFileBackupPage },
    { path: "restore/file", component: DecryptFileBackupPage },
    { path: ".*", component: NotFoundPageState },
  ];
</script>

<main>
  <Router {routes} />
</main>
