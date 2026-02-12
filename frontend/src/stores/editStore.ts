import { create } from 'zustand';
import type { QueueJob } from '../lib/types';

interface EditStoreState {
  jobs: QueueJob[];
  queueLocked: boolean;
  addJobs: (jobs: QueueJob[]) => void;
  removeJob: (jobID: string) => void;
  clearJobs: () => void;
  reorderJobs: (from: number, to: number) => void;
  setQueueLocked: (locked: boolean) => void;
}

export const useEditStore = create<EditStoreState>((set, get) => ({
  jobs: [],
  queueLocked: false,
  addJobs: (jobs) =>
    set((state) => ({
      jobs: [...state.jobs, ...jobs]
    })),
  removeJob: (jobID) =>
    set((state) => ({
      jobs: state.jobs.filter((j) => j.job_id !== jobID)
    })),
  clearJobs: () => set({ jobs: [] }),
  reorderJobs: (from, to) => {
    const jobs = [...get().jobs];
    const [item] = jobs.splice(from, 1);
    jobs.splice(to, 0, item);
    set({ jobs });
  },
  setQueueLocked: (locked) => set({ queueLocked: locked })
}));
