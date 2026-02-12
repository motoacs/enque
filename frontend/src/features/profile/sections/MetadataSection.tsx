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
      <h3 className="section-heading mb-3">
        {t("profile.metadata")}
      </h3>
      <div className="space-y-2">
        <label className="flex items-center gap-2 text-xs font-medium cursor-pointer" style={{ color: '#e8e6e3' }}>
          <input
            type="checkbox"
            checked={allOn}
            onChange={toggleAll}
            disabled={isPreset}
            className="rounded"
          />
          {t("profile.metadataAll")}
        </label>
        <div className="ml-5 space-y-1.5">
          {toggles.map(({ label, field }) => (
            <label key={field} className="flex items-center gap-2 text-xs cursor-pointer" style={{ color: '#9d9da7' }}>
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
