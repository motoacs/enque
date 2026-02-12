import { useTranslation } from "react-i18next";
import { useProfileStore } from "@/stores/profileStore";

export function MetadataSection() {
  const { t } = useTranslation();
  const { editingProfile: p, updateEditingProfile: update } = useProfileStore();
  if (!p) return null;

  const isPreset = p.is_preset;

  const toggles = [
    { label: "Metadata Copy", field: "metadata_copy" as const },
    { label: "Video Metadata", field: "video_metadata_copy" as const },
    { label: "Audio Metadata", field: "audio_metadata_copy" as const },
    { label: "Chapter Copy", field: "chapter_copy" as const },
    { label: "Subtitle Copy", field: "sub_copy" as const },
    { label: "Data Copy", field: "data_copy" as const },
    { label: "Attachment Copy", field: "attachment_copy" as const },
    { label: "File Time Restore", field: "restore_file_time" as const },
  ];

  const allOn = toggles.every((t) => p[t.field]);

  const toggleAll = () => {
    if (isPreset) return;
    const newValue = !allOn;
    const updates: Record<string, boolean> = {};
    toggles.forEach((t) => { updates[t.field] = newValue; });
    update(updates as any);
  };

  return (
    <section>
      <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">
        {t("profile.metadata")}
      </h3>
      <div className="space-y-1.5">
        <label className="flex items-center gap-1.5 text-xs text-zinc-300 font-medium">
          <input
            type="checkbox"
            checked={allOn}
            onChange={toggleAll}
            disabled={isPreset}
            className="rounded"
          />
          {t("profile.metadataAll")}
        </label>
        <div className="ml-4 space-y-1">
          {toggles.map(({ label, field }) => (
            <label key={field} className="flex items-center gap-1.5 text-xs text-zinc-400">
              <input
                type="checkbox"
                checked={p[field]}
                onChange={(e) => update({ [field]: e.target.checked } as any)}
                disabled={isPreset}
                className="rounded"
              />
              {label}
            </label>
          ))}
        </div>
      </div>
    </section>
  );
}
