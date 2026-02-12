import { create } from "zustand";

export interface ToolInfo {
  name: string;
  path: string;
  version: string;
  found: boolean;
  error?: string;
  supported: boolean;
}

export interface AppConfig {
  version: number;
  nvencc_path: string;
  qsvenc_path: string;
  ffmpeg_path: string;
  ffprobe_path: string;
  max_concurrent_jobs: number;
  on_error: string;
  decoder_fallback: boolean;
  keep_failed_temp: boolean;
  no_output_timeout_sec: number;
  no_progress_timeout_sec: number;
  post_complete_action: string;
  post_complete_command: string;
  output_folder_mode: string;
  output_folder_path: string;
  output_name_template: string;
  output_container: string;
  overwrite_mode: string;
  language: string;
  default_profile_id: string;
}

export interface DetectionResult {
  nvencc: ToolInfo;
  qsvenc: ToolInfo;
  ffmpeg: ToolInfo;
  ffprobe: ToolInfo;
}

interface AppState {
  config: AppConfig | null;
  tools: DetectionResult | null;
  gpuInfo: string;
  language: string;
  setConfig: (config: AppConfig) => void;
  setTools: (tools: DetectionResult) => void;
  setGPUInfo: (info: string) => void;
  setLanguage: (lang: string) => void;
  updateConfig: (updates: Partial<AppConfig>) => void;
}

export const useAppStore = create<AppState>((set) => ({
  config: null,
  tools: null,
  gpuInfo: "",
  language: "ja",
  setConfig: (config) => set({ config, language: config.language }),
  setTools: (tools) => set({ tools }),
  setGPUInfo: (gpuInfo) => set({ gpuInfo }),
  setLanguage: (language) => set({ language }),
  updateConfig: (updates) =>
    set((s) => ({
      config: s.config ? { ...s.config, ...updates } : null,
    })),
}));
