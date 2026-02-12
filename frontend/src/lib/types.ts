export type EncoderType = 'nvencc' | 'qsvenc' | 'ffmpeg';
export type JobStatus = 'pending' | 'running' | 'completed' | 'failed' | 'cancelled' | 'timeout' | 'skipped';

export interface JobProgress {
  percent: number | null;
  fps: number | null;
  bitrate_kbps: number | null;
  eta_sec: number | null;
  raw_line?: string;
}

export interface QueueJob {
  job_id: string;
  input_path: string;
  input_size_bytes: number;
  status: JobStatus;
  progress: JobProgress;
  started_at?: string;
  finished_at?: string;
  worker_id?: number;
  exit_code?: number;
  error_message: string;
}

export interface NVEncCAdvanced {
  interlace: string;
  avsw_decoder: string;
  input_csp: string;
  output_csp: string;
  tune: string;
  max_bitrate: number | null;
  vbr_quality: number | null;
  lookahead_level: number | null;
  weightp: boolean;
  mv_precision: string;
  refs_forward: number | null;
  refs_backward: number | null;
  level: string;
  profile: string;
  tier: string;
  ssim: boolean;
  psnr: boolean;
  trim: string;
  seek: string;
  seekto: string;
  video_metadata: string;
  audio_copy: string;
  audio_codec: string;
  audio_bitrate: string;
  audio_quality: string;
  audio_samplerate: string;
  audio_metadata: string;
  sub_copy: string;
  sub_metadata: string;
  data_copy: string;
  attachment_copy: string;
  metadata: string;
  output_thread: number | null;
}

export interface Profile {
  id: string;
  version: number;
  name: string;
  is_preset: boolean;
  encoder_type: EncoderType;
  encoder_options: Record<string, unknown>;
  codec: 'h264' | 'hevc' | 'av1';
  rate_control: 'qvbr' | 'cqp' | 'cbr' | 'vbr';
  rate_value: number;
  preset: 'P1' | 'P2' | 'P3' | 'P4' | 'P5' | 'P6' | 'P7';
  output_depth: 8 | 10;
  multipass: 'none' | 'quarter' | 'full';
  output_res: string;
  bframes: number | null;
  ref: number | null;
  lookahead: number | null;
  gop_len: number | null;
  aq: boolean;
  aq_temporal: boolean;
  split_enc: 'off' | 'auto' | 'auto_forced' | 'forced_2' | 'forced_3' | 'forced_4';
  parallel: 'off' | 'auto' | '2' | '3';
  decoder: 'avhw' | 'avsw';
  device: string;
  audio_mode: 'copy' | 'aac' | 'opus';
  audio_bitrate: number;
  colormatrix: string;
  transfer: string;
  colorprim: string;
  colorrange: string;
  dhdr10_info: string;
  metadata_copy: boolean;
  video_metadata_copy: boolean;
  audio_metadata_copy: boolean;
  chapter_copy: boolean;
  sub_copy: boolean;
  data_copy: boolean;
  attachment_copy: boolean;
  restore_file_time: boolean;
  nvencc_advanced: NVEncCAdvanced;
  custom_options: string;
}

export interface AppConfig {
  version: number;
  nvencc_path: string;
  qsvenc_path: string;
  ffmpeg_path: string;
  ffprobe_path: string;
  max_concurrent_jobs: number;
  on_error: 'skip' | 'stop';
  decoder_fallback: boolean;
  keep_failed_temp: boolean;
  no_output_timeout_sec: number;
  no_progress_timeout_sec: number;
  post_complete_action: 'none' | 'shutdown' | 'sleep' | 'custom';
  post_complete_command: string;
  output_folder_mode: 'same_as_input' | 'specified';
  output_folder_path: string;
  output_name_template: string;
  output_container: string;
  overwrite_mode: 'ask' | 'auto_rename';
  language: 'ja' | 'en';
  default_profile_id: string;
}

export interface ToolInfo {
  found: boolean;
  path: string;
  version: string;
  warning?: string;
}

export interface ToolSnapshot {
  nvencc: ToolInfo;
  qsvenc: ToolInfo;
  ffmpeg: ToolInfo;
  ffprobe: ToolInfo;
}

export interface GPUInfo {
  check_device_output: string;
  check_features_output: string;
}

export interface BootstrapResponse {
  config: AppConfig;
  profiles: Profile[];
  tools: ToolSnapshot;
  gpu_info: GPUInfo;
  warnings: string[];
}

export interface PreviewCommandResponse {
  argv: string[];
  display_command: string;
}
