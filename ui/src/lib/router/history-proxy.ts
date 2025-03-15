let proxied = false;

export function proxyHistory() {
  if (proxied) {
    return;
  }

  window.history.pushState = new Proxy(window.history.pushState, {
    apply: (target, thisArg, argumentsList: [any, any, any]) => {
      const result = target.apply(thisArg, argumentsList);

      window.dispatchEvent(new Event("pushstate"));

      return result;
    },
  });

  window.history.replaceState = new Proxy(window.history.replaceState, {
    apply: (target, thisArg, argumentsList: [any, any, any]) => {
      const result = target.apply(thisArg, argumentsList);

      window.dispatchEvent(new Event("replacestate"));

      return result;
    },
  });

  proxied = true;
}
