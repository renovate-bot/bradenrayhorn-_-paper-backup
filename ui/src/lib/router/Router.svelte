<script lang="ts">
  import { onMount, type Component } from "svelte";
  import { proxyHistory } from "./history-proxy";

  let RouteComponent = $state<Component | null>(null);

  let lastUri: string | undefined = undefined;

  export type Route = {
    path: string;
    component: Component;
    orRedirect?: () => string | undefined;
  };
  const props: { routes: Route[]; basePath?: string } = $props();

  const routes = $derived.by(() => {
    return props.routes.map(({ path, ...route }) => {
      const basePath = (props.basePath ?? "").replace(/^\/+|\/+$/g, "");
      path = path.replace(/^\/+|\/+$/g, "");

      const regex = new RegExp(`${basePath}/${path.replaceAll("/", "\\/")}$`);

      return { path: regex, ...route };
    });
  });

  function updateRoute(uri: string) {
    uri = uri.replace("#", "");
    if (!uri.startsWith("/")) {
      uri = "/" + uri;
    }

    if (uri === lastUri) {
      return;
    }
    lastUri = uri;

    let route: Component | null = null;
    const now = performance.now();
    for (const candidate of routes) {
      if (candidate.path.test(uri)) {
        if (candidate.orRedirect) {
          const goTo = candidate.orRedirect();
          if (goTo) {
            window.history.replaceState(undefined, "", goTo);
            updateRoute(goTo);
            return;
          }
        }
        route = candidate.component;
        break;
      }
    }
    console.log(`router match in ${performance.now() - now}ms`, uri);

    RouteComponent = route;
  }

  onMount(() => {
    const doUpdate = () => updateRoute(location.hash);

    // pushstate and replacestate are custom events created by proxyHistory()
    window.addEventListener("popstate", doUpdate);
    window.addEventListener("pushstate", doUpdate);
    window.addEventListener("replacestate", doUpdate);

    return () => {
      window.removeEventListener("popstate", doUpdate);
      window.removeEventListener("pushstate", doUpdate);
      window.removeEventListener("replacestate", doUpdate);
    };
  });

  proxyHistory();
  updateRoute(location.hash);
</script>

<RouteComponent />
