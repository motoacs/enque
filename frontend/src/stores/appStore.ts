import { create } from 'zustand';
import type { AppConfig, GPUInfo, ToolSnapshot } from '../lib/types';

const defaultConfig: AppConfig = {
  version: 1,
  nvencc_path: '',
  qsvenc_path: '',
  ffmpeg_path: '',
  ffprobe_path: '',
  max_concurrent_jobs: 1,
  on_error: 'skip',
  decoder_fallback: false,
  keep_failed_temp: false,
  no_output_timeout_sec: 600,
  no_progress_timeout_sec: 300,
  post_complete_action: 'none',
  post_complete_command: '',
  output_folder_mode: 'same_as_input',
  output_folder_path: '',
  output_name_template: '{name}_encoded.{ext}',
  output_container: 'mkv',
  overwrite_mode: 'ask',
  language: 'ja',
  default_profile_id: ''
};

const emptyTools: ToolSnapshot = {
  nvencc: { found: false, path: '', version: '' },
  qsvenc: { found: false, path: '', version: '' },
  ffmpeg: { found: false, path: '', version: '' },
  ffprobe: { found: false, path: '', version: '' }
};

const emptyGPU: GPUInfo = {
  check_device_output: '',
  check_features_output: ''
};

interface AppStoreState {
  config: AppConfig;
  tools: ToolSnapshot;
  gpuInfo: GPUInfo;
  warnings: string[];
  setConfig: (next: AppConfig) => void;
  patchConfig: (patch: Partial<AppConfig>) => void;
  setTools: (tools: ToolSnapshot) => void;
  setGPUInfo: (gpuInfo: GPUInfo) => void;
  setWarnings: (warnings: string[]) => void;
}

export const useAppStore = create<AppStoreState>((set) => ({
  config: defaultConfig,
  tools: emptyTools,
  gpuInfo: emptyGPU,
  warnings: [],
  setConfig: (next) => set({ config: next }),
  patchConfig: (patch) => set((state) => ({ config: { ...state.config, ...patch } })),
  setTools: (tools) => set({ tools }),
  setGPUInfo: (gpuInfo) => set({ gpuInfo }),
  setWarnings: (warnings) => set({ warnings })
}));
