import { create } from "zustand";

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
  encoder_type: string;
  encoder_options: Record<string, unknown>;
  codec: string;
  rate_control: string;
  rate_value: number;
  preset: string;
  output_depth: number;
  multipass: string;
  output_res: string;
  bframes: number | null;
  ref: number | null;
  lookahead: number | null;
  gop_len: number | null;
  aq: boolean;
  aq_temporal: boolean;
  split_enc: string;
  parallel: string;
  decoder: string;
  device: string;
  audio_mode: string;
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

interface ProfileState {
  profiles: Profile[];
  editingProfile: Profile | null;
  setProfiles: (profiles: Profile[]) => void;
  setEditingProfile: (profile: Profile | null) => void;
  updateEditingProfile: (updates: Partial<Profile>) => void;
  updateAdvanced: (updates: Partial<NVEncCAdvanced>) => void;
}

export const useProfileStore = create<ProfileState>((set) => ({
  profiles: [],
  editingProfile: null,
  setProfiles: (profiles) => set({ profiles }),
  setEditingProfile: (profile) => set({ editingProfile: profile ? { ...profile } : null }),
  updateEditingProfile: (updates) =>
    set((s) => ({
      editingProfile: s.editingProfile
        ? { ...s.editingProfile, ...updates }
        : null,
    })),
  updateAdvanced: (updates) =>
    set((s) => ({
      editingProfile: s.editingProfile
        ? {
            ...s.editingProfile,
            nvencc_advanced: { ...s.editingProfile.nvencc_advanced, ...updates },
          }
        : null,
    })),
}));
