import { create } from 'zustand';
import type { JobStatus } from '../lib/types';

interface JobRuntime {
  status: JobStatus;
  percent: number | null;
  fps: number | null;
  bitrate_kbps: number | null;
  eta_sec: number | null;
  logs: string[];
  final_output_path?: string;
  exit_code?: number;
  error_message?: string;
}

interface EncodeStoreState {
  sessionId: string;
  state: 'idle' | 'running' | 'stopping' | 'aborting' | 'completed' | 'aborted';
  totalJobs: number;
  completedJobs: number;
  jobs: Record<string, JobRuntime>;
  upsertJob: (jobID: string, patch: Partial<JobRuntime>) => void;
  appendJobLog: (jobID: string, line: string) => void;
  setSession: (sessionId: string, totalJobs: number) => void;
  setState: (state: EncodeStoreState['state']) => void;
  setCompletedJobs: (n: number) => void;
  reset: () => void;
}

const LOG_LIMIT = 2000;

export const useEncodeStore = create<EncodeStoreState>((set) => ({
  sessionId: '',
  state: 'idle',
  totalJobs: 0,
  completedJobs: 0,
  jobs: {},
  upsertJob: (jobID, patch) =>
    set((state) => {
      const prev = state.jobs[jobID] ?? {
        status: 'pending' as const,
        percent: null,
        fps: null,
        bitrate_kbps: null,
        eta_sec: null,
        logs: []
      };
      return {
        jobs: {
          ...state.jobs,
          [jobID]: {
            ...prev,
            ...patch
          }
        }
      };
    }),
  appendJobLog: (jobID, line) =>
    set((state) => {
      const prev = state.jobs[jobID] ?? {
        status: 'pending' as const,
        percent: null,
        fps: null,
        bitrate_kbps: null,
        eta_sec: null,
        logs: []
      };
      const logs = [...prev.logs, line];
      if (logs.length > LOG_LIMIT) {
        logs.splice(0, logs.length - LOG_LIMIT);
      }
      return {
        jobs: {
          ...state.jobs,
          [jobID]: {
            ...prev,
            logs
          }
        }
      };
    }),
  setSession: (sessionId, totalJobs) =>
    set({ sessionId, totalJobs, state: 'running', completedJobs: 0, jobs: {} }),
  setState: (state) => set({ state }),
  setCompletedJobs: (n) => set({ completedJobs: n }),
  reset: () =>
    set({
      sessionId: '',
      state: 'idle',
      totalJobs: 0,
      completedJobs: 0,
      jobs: {}
    })
}));
