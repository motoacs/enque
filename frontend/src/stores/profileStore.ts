import { create } from 'zustand';
import type { Profile } from '../lib/types';

interface ProfileStoreState {
  profiles: Profile[];
  selectedProfileId: string;
  validationErrors: Record<string, string>;
  setProfiles: (profiles: Profile[]) => void;
  selectProfile: (profileId: string) => void;
  updateSelectedProfile: (patch: Partial<Profile>) => void;
  setValidationErrors: (errors: Record<string, string>) => void;
}

export const useProfileStore = create<ProfileStoreState>((set, get) => ({
  profiles: [],
  selectedProfileId: '',
  validationErrors: {},
  setProfiles: (profiles) =>
    set((state) => ({
      profiles,
      selectedProfileId: state.selectedProfileId || profiles[0]?.id || ''
    })),
  selectProfile: (profileId) => set({ selectedProfileId: profileId }),
  updateSelectedProfile: (patch) => {
    const { profiles, selectedProfileId } = get();
    const next = profiles.map((p) => (p.id === selectedProfileId ? { ...p, ...patch } : p));
    set({ profiles: next });
  },
  setValidationErrors: (errors) => set({ validationErrors: errors })
}));

export function useSelectedProfile(): Profile | undefined {
  const { profiles, selectedProfileId } = useProfileStore.getState();
  return profiles.find((p) => p.id === selectedProfileId);
}
