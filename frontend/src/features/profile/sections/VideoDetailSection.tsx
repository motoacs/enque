import { useTranslation } from "react-i18next";
import { useProfileStore } from "@/stores/profileStore";

export function VideoDetailSection() {
  const { t } = useTranslation();
  const { editingProfile: p, updateEditingProfile: update } = useProfileStore();
  if (!p) return null;

  const isPreset = p.is_preset;

  const setNullableInt = (field: string, value: string) => {
    update({ [field]: value === "" ? null : Number(value) } as any);
  };

  return (
    <section>
      <h3 className="section-heading mb-3">
        {t("profile.videoDetail")}
      </h3>
      <div className="space-y-2.5">
        {[
          { label: "B-Frames", field: "bframes", min: 0, max: 7 },
          { label: "Ref Frames", field: "ref", min: 0, max: 16 },
          { label: "Lookahead", field: "lookahead", min: 0, max: 32 },
          { label: "GOP Length", field: "gop_len", min: 0, max: 9999 },
        ].map(({ label, field, min, max }) => (
          <div key={field} className="flex items-center gap-2">
            <label className="form-label w-28">{label}</label>
            <input
              type="number"
              value={(p as any)[field] ?? ""}
              onChange={(e) => setNullableInt(field, e.target.value)}
              disabled={isPreset}
              min={min}
              max={max}
              placeholder="auto"
              className="w-20 form-input font-mono"
            />
          </div>
        ))}

        <div className="flex items-center gap-5 pt-1">
          <label className="flex items-center gap-2 text-xs cursor-pointer" style={{ color: '#9d9da7' }}>
            <input
              type="checkbox"
              checked={p.aq}
              onChange={(e) => update({ aq: e.target.checked })}
              disabled={isPreset}
              className="rounded"
            />
            AQ
          </label>
          <label className="flex items-center gap-2 text-xs cursor-pointer" style={{ color: '#9d9da7' }}>
            <input
              type="checkbox"
              checked={p.aq_temporal}
              onChange={(e) => update({ aq_temporal: e.target.checked })}
              disabled={isPreset}
              className="rounded"
            />
            AQ Temporal
          </label>
        </div>
      </div>
    </section>
  );
}
