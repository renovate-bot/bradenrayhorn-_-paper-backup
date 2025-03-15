export function link(node: HTMLElement, fn?: () => void) {
  function onClick(e: MouseEvent) {
    if (fn) {
      fn();
    }
    e.preventDefault();
    window.history.pushState(undefined, "", node.getAttribute("href"));
  }

  $effect(() => {
    node.addEventListener("click", onClick);

    return () => {
      node.removeEventListener("click", onClick);
    };
  });
}
