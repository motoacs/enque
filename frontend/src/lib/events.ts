// Wails event listener types and registration.

import { useEncodeStore } from "@/stores/encodeStore";

export type EventCallback = (data: unknown) => void;

export const EventNames = {
  SESSION_STARTED: "enque:session_started",
  JOB_STARTED: "enque:job_started",
  JOB_PROGRESS: "enque:job_progress",
  JOB_LOG: "enque:job_log",
  JOB_NEEDS_OVERWRITE: "enque:job_needs_overwrite",
  JOB_FINISHED: "enque:job_finished",
  SESSION_STATE: "enque:session_state",
  SESSION_FINISHED: "enque:session_finished",
  WARNING: "enque:warning",
  ERROR: "enque:error",
} as const;

// Register all Wails event listeners. Call once on app startup.
export function registerEventListeners() {
  const runtime = (window as unknown as Record<string, unknown>).runtime;
  if (!runtime) {
    console.warn("Wails runtime not available, event listeners not registered");
    return;
  }

  const eventsOn = (runtime as Record<string, Function>).EventsOn;
  if (!eventsOn) return;

  const store = useEncodeStore.getState;

  eventsOn(EventNames.SESSION_STARTED, (data: Record<string, unknown>) => {
    store().onSessionStarted(data);
  });

  eventsOn(EventNames.JOB_STARTED, (data: Record<string, unknown>) => {
    store().onJobStarted(data);
  });

  eventsOn(EventNames.JOB_PROGRESS, (data: Record<string, unknown>) => {
    store().onJobProgress(data);
  });

  eventsOn(EventNames.JOB_LOG, (data: Record<string, unknown>) => {
    store().onJobLog(data);
  });

  eventsOn(EventNames.JOB_NEEDS_OVERWRITE, (data: Record<string, unknown>) => {
    store().onJobNeedsOverwrite(data);
  });

  eventsOn(EventNames.JOB_FINISHED, (data: Record<string, unknown>) => {
    store().onJobFinished(data);
  });

  eventsOn(EventNames.SESSION_STATE, (data: Record<string, unknown>) => {
    store().onSessionState(data);
  });

  eventsOn(EventNames.SESSION_FINISHED, (data: Record<string, unknown>) => {
    store().onSessionFinished(data);
  });

  eventsOn(EventNames.WARNING, (data: Record<string, unknown>) => {
    store().onWarning(data);
  });

  eventsOn(EventNames.ERROR, (data: Record<string, unknown>) => {
    console.error("Enque error:", data);
  });
}
