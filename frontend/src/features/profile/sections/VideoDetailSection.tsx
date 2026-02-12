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
      <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">
        {t("profile.videoDetail")}
      </h3>
      <div className="space-y-2">
        {[
          { label: "B-Frames", field: "bframes", min: 0, max: 7 },
          { label: "Ref Frames", field: "ref", min: 0, max: 16 },
          { label: "Lookahead", field: "lookahead", min: 0, max: 32 },
          { label: "GOP Length", field: "gop_len", min: 0, max: 9999 },
        ].map(({ label, field, min, max }) => (
          <div key={field} className="flex items-center gap-2">
            <label className="text-xs text-zinc-400 w-28 shrink-0">{label}</label>
            <input
              type="number"
              value={(p as any)[field] ?? ""}
              onChange={(e) => setNullableInt(field, e.target.value)}
              disabled={isPreset}
              min={min}
              max={max}
              placeholder="auto"
              className="w-20 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60 placeholder:text-zinc-600"
            />
          </div>
        ))}

        <div className="flex items-center gap-4">
          <label className="flex items-center gap-1.5 text-xs text-zinc-400">
            <input
              type="checkbox"
              checked={p.aq}
              onChange={(e) => update({ aq: e.target.checked })}
              disabled={isPreset}
              className="rounded"
            />
            AQ
          </label>
          <label className="flex items-center gap-1.5 text-xs text-zinc-400">
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
