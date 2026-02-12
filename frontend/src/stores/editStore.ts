import { create } from "zustand";

export interface QueueJob {
  jobId: string;
  inputPath: string;
  inputSizeBytes: number;
  fileName: string;
}

export interface OutputSettings {
  outputFolderMode: string;
  outputFolderPath: string;
  outputNameTemplate: string;
  outputContainer: string;
  overwriteMode: string;
}

interface EditState {
  jobs: QueueJob[];
  selectedProfileId: string;
  outputSettings: OutputSettings;
  addJobs: (jobs: QueueJob[]) => void;
  removeJob: (jobId: string) => void;
  reorderJobs: (jobs: QueueJob[]) => void;
  clearJobs: () => void;
  setSelectedProfileId: (id: string) => void;
  setOutputSettings: (settings: Partial<OutputSettings>) => void;
}

export const useEditStore = create<EditState>((set) => ({
  jobs: [],
  selectedProfileId: "",
  outputSettings: {
    outputFolderMode: "same_as_input",
    outputFolderPath: "",
    outputNameTemplate: "{name}_encoded.{ext}",
    outputContainer: "mkv",
    overwriteMode: "ask",
  },
  addJobs: (newJobs) => set((s) => ({ jobs: [...s.jobs, ...newJobs] })),
  removeJob: (jobId) =>
    set((s) => ({ jobs: s.jobs.filter((j) => j.jobId !== jobId) })),
  reorderJobs: (jobs) => set({ jobs }),
  clearJobs: () => set({ jobs: [] }),
  setSelectedProfileId: (id) => set({ selectedProfileId: id }),
  setOutputSettings: (settings) =>
    set((s) => ({
      outputSettings: { ...s.outputSettings, ...settings },
    })),
}));
