import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useProfileStore, type Profile } from "@/stores/profileStore";
import { useEditStore } from "@/stores/editStore";
import { Plus, Copy, Pencil, Trash2, Lock, Layers } from "lucide-react";
import { VideoBasicSection } from "./sections/VideoBasicSection";
import { VideoDetailSection } from "./sections/VideoDetailSection";
import { SpeedSection } from "./sections/SpeedSection";
import { AudioSection } from "./sections/AudioSection";
import { ColorSection } from "./sections/ColorSection";
import { MetadataSection } from "./sections/MetadataSection";
import { AdvancedSection } from "./sections/AdvancedSection";
import { CustomOptionsSection } from "./sections/CustomOptionsSection";
import { ProfileNameDialog } from "./ProfileNameDialog";
import { DeleteConfirmDialog } from "./DeleteConfirmDialog";
import * as api from "@/lib/api";

type NameDialogMode = "create" | "duplicate" | "rename" | null;

export function ProfileEditor() {
  const { t } = useTranslation();
  const {
    profiles,
    editingProfile,
    setEditingProfile,
    addProfile,
    removeProfile,
    updateProfileInList,
  } = useProfileStore();
  const { selectedProfileId, setSelectedProfileId } = useEditStore();

  const [nameDialogMode, setNameDialogMode] = useState<NameDialogMode>(null);
  const [nameDialogInitial, setNameDialogInitial] = useState("");
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const builtinProfiles = profiles.filter((p) => p.is_preset);
  const userProfiles = profiles.filter((p) => !p.is_preset);
  const selectedProfile = profiles.find((p) => p.id === selectedProfileId);
  const isPreset = selectedProfile?.is_preset ?? true;

  const handleProfileSelect = (id: string) => {
    setSelectedProfileId(id);
    const p = profiles.find((p) => p.id === id);
    if (p) setEditingProfile(p);
  };

  const selectProfile = (profile: Profile) => {
    setSelectedProfileId(profile.id);
    setEditingProfile(profile);
  };

  // --- Name dialog helpers ---

  const getNameDialogTitle = () => {
    switch (nameDialogMode) {
      case "create": return t("profile.create");
      case "duplicate": return t("profile.duplicate");
      case "rename": return t("profile.rename");
      default: return "";
    }
  };

  const handleNameConfirm = async (name: string) => {
    try {
      if (nameDialogMode === "create") {
        const base = builtinProfiles[0];
        if (!base) return;
        const newProfile = await api.duplicateProfile(base.id, name) as Profile;
        addProfile(newProfile);
        selectProfile(newProfile);
      } else if (nameDialogMode === "duplicate") {
        if (!selectedProfile) return;
        const newProfile = await api.duplicateProfile(selectedProfile.id, name) as Profile;
        addProfile(newProfile);
        selectProfile(newProfile);
      } else if (nameDialogMode === "rename") {
        if (!editingProfile) return;
        const updated = { ...editingProfile, name };
        await api.upsertProfile(updated);
        updateProfileInList(updated as Profile);
        setEditingProfile(updated as Profile);
      }
    } catch (err) {
      console.error("Profile operation failed:", err);
    }
    setNameDialogMode(null);
  };

  // --- Action handlers ---

  const handleCreateNew = () => {
    setNameDialogInitial(t("profile.newDefault"));
    setNameDialogMode("create");
  };

  const handleDuplicate = () => {
    if (!selectedProfile) return;
    setNameDialogInitial(selectedProfile.name + " (Copy)");
    setNameDialogMode("duplicate");
  };

  const handleRename = () => {
    if (!selectedProfile || selectedProfile.is_preset) return;
    setNameDialogInitial(selectedProfile.name);
    setNameDialogMode("rename");
  };

  const handleDeleteClick = () => {
    if (!selectedProfile || selectedProfile.is_preset) return;
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (!selectedProfile) return;
    try {
      await api.deleteProfile(selectedProfile.id);
      removeProfile(selectedProfile.id);
      // Select the first remaining profile
      const remaining = profiles.filter((p) => p.id !== selectedProfile.id);
      if (remaining.length > 0) {
        selectProfile(remaining[0]);
      } else {
        setSelectedProfileId("");
        setEditingProfile(null);
      }
    } catch (err) {
      console.error("Delete failed:", err);
    }
    setDeleteDialogOpen(false);
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center gap-2 px-4 py-2.5" style={{ borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
        <h2 className="font-display text-xs font-semibold uppercase tracking-wider shrink-0" style={{ color: '#9d9da7' }}>
          {t("profile.title")}
        </h2>
        <select
          value={selectedProfileId}
          onChange={(e) => handleProfileSelect(e.target.value)}
          className="flex-1 form-input"
        >
          <option value="">{t("profile.select")}</option>
          {builtinProfiles.length > 0 && (
            <optgroup label={t("profile.builtinGroup")}>
              {builtinProfiles.map((p) => (
                <option key={p.id} value={p.id}>{p.name}</option>
              ))}
            </optgroup>
          )}
          {userProfiles.length > 0 && (
            <optgroup label={t("profile.userGroup")}>
              {userProfiles.map((p) => (
                <option key={p.id} value={p.id}>{p.name}</option>
              ))}
            </optgroup>
          )}
        </select>
        <button onClick={handleCreateNew} className="icon-btn" title={t("profile.create")}>
          <Plus size={14} />
        </button>
        {selectedProfile && (
          <button onClick={handleDuplicate} className="icon-btn" title={t("profile.duplicate")}>
            <Copy size={14} />
          </button>
        )}
        {selectedProfile && !isPreset && (
          <>
            <button onClick={handleRename} className="icon-btn" title={t("profile.rename")}>
              <Pencil size={14} />
            </button>
            <button onClick={handleDeleteClick} className="icon-btn" title={t("profile.delete")}>
              <Trash2 size={14} />
            </button>
          </>
        )}
      </div>

      {editingProfile ? (
        <div className="flex-1 overflow-y-auto px-4 py-4 space-y-5">
          {isPreset && (
            <div
              className="flex items-center gap-2 rounded-lg px-3 py-2 text-xs"
              style={{
                background: 'rgba(255,255,255,0.03)',
                border: '1px solid rgba(255,255,255,0.06)',
                color: '#9d9da7',
              }}
            >
              <Lock size={12} />
              {t("profile.presetReadonly")}
            </div>
          )}
          <VideoBasicSection />
          <VideoDetailSection />
          <SpeedSection />
          <AudioSection />
          <ColorSection />
          <MetadataSection />
          <AdvancedSection />
          <CustomOptionsSection />
        </div>
      ) : (
        <div className="flex-1 flex flex-col items-center justify-center gap-3">
          <Layers size={24} style={{ color: '#3c3c48' }} />
          <span className="text-xs" style={{ color: '#5c5c68' }}>{t("profile.select")}</span>
        </div>
      )}

      <ProfileNameDialog
        open={nameDialogMode !== null}
        title={getNameDialogTitle()}
        initialName={nameDialogInitial}
        onConfirm={handleNameConfirm}
        onCancel={() => setNameDialogMode(null)}
      />

      <DeleteConfirmDialog
        open={deleteDialogOpen}
        profileName={selectedProfile?.name ?? ""}
        onConfirm={handleDeleteConfirm}
        onCancel={() => setDeleteDialogOpen(false)}
      />
    </div>
  );
}
