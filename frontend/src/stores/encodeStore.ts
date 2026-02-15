import { create } from "zustand";

export type SessionState = "idle" | "running" | "stopping" | "aborting" | "completed" | "aborted";

export interface JobProgress {
  jobId: string;
  percent?: number;
  fps?: number;
  bitrateKbps?: number;
  etaSec?: number;
  status: string;
  workerID?: number;
  inputPath?: string;
  inputSizeBytes?: number;
  tempOutputPath?: string;
  finalOutputPath?: string;
  exitCode?: number;
  errorMessage?: string;
}

export interface SessionSummary {
  sessionId: string;
  state: string;
  totalJobs: number;
  completedJobs: number;
  failedJobs: number;
  cancelledJobs: number;
  timeoutJobs: number;
  skippedJobs: number;
}

export interface OverwriteRequest {
  sessionId: string;
  jobId: string;
  outputPath: string;
}

const MAX_LOG_LINES = 2000;

interface EncodeState {
  sessionId: string;
  sessionState: SessionState;
  jobProgress: Record<string, JobProgress>;
  jobLogs: Record<string, string[]>;
  sessionSummary: SessionSummary | null;
  overwriteRequest: OverwriteRequest | null;
  warnings: string[];

  // Actions
  setSessionState: (state: SessionState) => void;
  initPendingJobs: (jobs: { jobId: string; inputPath: string }[]) => void;
  onSessionStarted: (data: Record<string, unknown>) => void;
  onJobStarted: (data: Record<string, unknown>) => void;
  onJobProgress: (data: Record<string, unknown>) => void;
  onJobLog: (data: Record<string, unknown>) => void;
  onJobFinished: (data: Record<string, unknown>) => void;
  onSessionState: (data: Record<string, unknown>) => void;
  onSessionFinished: (data: Record<string, unknown>) => void;
  onJobNeedsOverwrite: (data: Record<string, unknown>) => void;
  onWarning: (data: Record<string, unknown>) => void;
  skipPendingJob: (jobId: string) => void;
  clearOverwriteRequest: () => void;
  resetSession: () => void;
}

export const useEncodeStore = create<EncodeState>((set) => ({
  sessionId: "",
  sessionState: "idle",
  jobProgress: {},
  jobLogs: {},
  sessionSummary: null,
  overwriteRequest: null,
  warnings: [],

  setSessionState: (sessionState) => set({ sessionState }),

  initPendingJobs: (jobs) =>
    set((s) => {
      const progress: Record<string, JobProgress> = {};
      for (const j of jobs) {
        progress[j.jobId] = { jobId: j.jobId, status: "pending", inputPath: j.inputPath };
      }
      return { jobProgress: progress };
    }),

  onSessionStarted: (data) =>
    set((s) => ({
      sessionId: data.session_id as string,
      sessionState: "running",
      jobProgress: s.jobProgress,
      jobLogs: {},
      sessionSummary: null,
      overwriteRequest: null,
      warnings: [],
    })),

  onJobStarted: (data) =>
    set((s) => ({
      jobProgress: {
        ...s.jobProgress,
        [data.job_id as string]: {
          jobId: data.job_id as string,
          status: "running",
          workerID: data.worker_id as number,
          inputPath: data.input_path as string,
          inputSizeBytes: data.input_size_bytes as number,
          tempOutputPath: data.temp_output_path as string,
          finalOutputPath: data.final_output_path as string,
        },
      },
    })),

  onJobProgress: (data) =>
    set((s) => {
      const jobId = data.job_id as string;
      const existing = s.jobProgress[jobId] || { jobId, status: "running" };
      return {
        jobProgress: {
          ...s.jobProgress,
          [jobId]: {
            ...existing,
            percent: data.percent as number | undefined,
            fps: data.fps as number | undefined,
            bitrateKbps: data.bitrate_kbps as number | undefined,
            etaSec: data.eta_sec as number | undefined,
          },
        },
      };
    }),

  onJobLog: (data) =>
    set((s) => {
      const jobId = data.job_id as string;
      const line = data.line as string;
      const existing = s.jobLogs[jobId] || [];
      const updated = [...existing, line];
      // Keep ring buffer at MAX_LOG_LINES
      if (updated.length > MAX_LOG_LINES) {
        updated.splice(0, updated.length - MAX_LOG_LINES);
      }
      return {
        jobLogs: { ...s.jobLogs, [jobId]: updated },
      };
    }),

  onJobFinished: (data) =>
    set((s) => {
      const jobId = data.job_id as string;
      const existing = s.jobProgress[jobId] || { jobId, status: "pending" };
      return {
        jobProgress: {
          ...s.jobProgress,
          [jobId]: {
            ...existing,
            status: data.status as string,
            exitCode: data.exit_code as number | undefined,
            errorMessage: data.error_message as string | undefined,
          },
        },
      };
    }),

  onSessionState: (data) => {
    const state = data.state as string;
    let sessionState: SessionState = "running";
    if (state === "stopping") sessionState = "stopping";
    else if (state === "aborting") sessionState = "aborting";
    else if (state === "completed") sessionState = "completed";
    else if (state === "aborted") sessionState = "aborted";
    set({ sessionState });
  },

  onSessionFinished: (data) =>
    set({
      sessionState: (data.state as string) === "aborted" ? "aborted" : "completed",
      sessionSummary: {
        sessionId: data.session_id as string,
        state: data.state as string,
        totalJobs: data.total_jobs as number,
        completedJobs: data.completed_jobs as number,
        failedJobs: data.failed_jobs as number,
        cancelledJobs: data.cancelled_jobs as number,
        timeoutJobs: data.timeout_jobs as number,
        skippedJobs: data.skipped_jobs as number,
      },
    }),

  onJobNeedsOverwrite: (data) =>
    set({
      overwriteRequest: {
        sessionId: data.session_id as string,
        jobId: data.job_id as string,
        outputPath: data.output_path as string,
      },
    }),

  onWarning: (data) =>
    set((s) => ({
      warnings: [...s.warnings, data.message as string].slice(-50),
    })),

  skipPendingJob: (jobId) =>
    set((s) => {
      const existing = s.jobProgress[jobId];
      if (!existing) return s;
      return {
        jobProgress: {
          ...s.jobProgress,
          [jobId]: { ...existing, status: "skipped" },
        },
      };
    }),

  clearOverwriteRequest: () => set({ overwriteRequest: null }),

  resetSession: () =>
    set({
      sessionId: "",
      sessionState: "idle",
      jobProgress: {},
      jobLogs: {},
      sessionSummary: null,
      overwriteRequest: null,
      warnings: [],
    }),
}));
