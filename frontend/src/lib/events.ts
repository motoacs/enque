export type EventHandler<T = unknown> = (payload: T) => void;

export function onEvent<T = unknown>(name: string, handler: EventHandler<T>): () => void {
  const runtime = window.runtime;
  if (!runtime?.EventsOn) {
    return () => {};
  }
  return runtime.EventsOn(name, handler as (payload: unknown) => void);
}
