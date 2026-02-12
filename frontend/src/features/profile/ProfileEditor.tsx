import { useTranslation } from "react-i18next";
import { useProfileStore, type Profile } from "@/stores/profileStore";
import { useEditStore } from "@/stores/editStore";
import { Copy } from "lucide-react";
import { VideoBasicSection } from "./sections/VideoBasicSection";
import { VideoDetailSection } from "./sections/VideoDetailSection";
import { SpeedSection } from "./sections/SpeedSection";
import { AudioSection } from "./sections/AudioSection";
import { ColorSection } from "./sections/ColorSection";
import { MetadataSection } from "./sections/MetadataSection";
import { AdvancedSection } from "./sections/AdvancedSection";
import { CustomOptionsSection } from "./sections/CustomOptionsSection";

export function ProfileEditor() {
  const { t } = useTranslation();
  const {
    profiles,
    editingProfile,
    setEditingProfile,
  } = useProfileStore();
  const { selectedProfileId, setSelectedProfileId } = useEditStore();

  const handleProfileSelect = (id: string) => {
    setSelectedProfileId(id);
    const p = profiles.find((p) => p.id === id);
    if (p) setEditingProfile(p);
  };

  const handleDuplicate = () => {
    if (!editingProfile) return;
    // In production, this would call the Wails API
  };

  const selectedProfile = profiles.find((p) => p.id === selectedProfileId);

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center gap-2 px-3 py-2 border-b border-zinc-700">
        <h2 className="text-sm font-medium text-zinc-300 shrink-0">
          {t("profile.title")}
        </h2>
        <select
          value={selectedProfileId}
          onChange={(e) => handleProfileSelect(e.target.value)}
          className="flex-1 bg-zinc-700 text-zinc-200 text-sm rounded px-2 py-1 border border-zinc-600 focus:outline-none focus:border-zinc-500"
        >
          <option value="">{t("profile.select")}</option>
          {profiles.map((p) => (
            <option key={p.id} value={p.id}>
              {p.is_preset ? `[${t("profile.preset")}] ` : ""}
              {p.name}
            </option>
          ))}
        </select>
        {selectedProfile?.is_preset && (
          <button
            onClick={handleDuplicate}
            className="p-1.5 rounded hover:bg-zinc-700 text-zinc-400 hover:text-zinc-200"
            title={t("profile.duplicate")}
          >
            <Copy size={14} />
          </button>
        )}
      </div>

      {editingProfile ? (
        <div className="flex-1 overflow-y-auto p-3 space-y-4">
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
        <div className="flex-1 flex items-center justify-center text-zinc-500 text-sm">
          {t("profile.select")}
        </div>
      )}
    </div>
  );
}
