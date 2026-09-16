// `main: import.meta.url` tells alchemy to bundle this file as the worker
// script (docs: infrastructure-as-effects/runtime).
interface ImportMeta {
  readonly url: string;
}
