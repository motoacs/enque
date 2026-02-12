import { create } from "zustand";
import * as api from "@/lib/api";

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
  output_container: string;
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
  addProfile: (profile: Profile) => void;
  removeProfile: (id: string) => void;
  updateProfileInList: (profile: Profile) => void;
}

// Debounced auto-save for user presets
let saveTimer: ReturnType<typeof setTimeout> | null = null;
let pendingSaveProfile: Profile | null = null;

function scheduleSave(profile: Profile) {
  pendingSaveProfile = profile;
  if (saveTimer) clearTimeout(saveTimer);
  saveTimer = setTimeout(() => {
    flushSave();
  }, 500);
}

function flushSave() {
  if (pendingSaveProfile) {
    api.upsertProfile(pendingSaveProfile).catch((err) => {
      console.error("Auto-save failed:", err);
    });
    pendingSaveProfile = null;
  }
  if (saveTimer) {
    clearTimeout(saveTimer);
    saveTimer = null;
  }
}

export const useProfileStore = create<ProfileState>((set, get) => ({
  profiles: [],
  editingProfile: null,
  setProfiles: (profiles) => set({ profiles }),
  setEditingProfile: (profile) => {
    // Flush any pending save for the previous profile before switching
    const prev = get().editingProfile;
    if (prev && !prev.is_preset && pendingSaveProfile) {
      flushSave();
    }
    set({ editingProfile: profile ? { ...profile } : null });
  },
  updateEditingProfile: (updates) =>
    set((s) => {
      if (!s.editingProfile) return {};
      const updated = { ...s.editingProfile, ...updates };
      if (!updated.is_preset) {
        scheduleSave(updated);
      }
      return { editingProfile: updated };
    }),
  updateAdvanced: (updates) =>
    set((s) => {
      if (!s.editingProfile) return {};
      const updated = {
        ...s.editingProfile,
        nvencc_advanced: { ...s.editingProfile.nvencc_advanced, ...updates },
      };
      if (!updated.is_preset) {
        scheduleSave(updated);
      }
      return { editingProfile: updated };
    }),
  addProfile: (profile) =>
    set((s) => ({ profiles: [...s.profiles, profile] })),
  removeProfile: (id) =>
    set((s) => ({ profiles: s.profiles.filter((p) => p.id !== id) })),
  updateProfileInList: (profile) =>
    set((s) => ({
      profiles: s.profiles.map((p) => (p.id === profile.id ? profile : p)),
    })),
}));
