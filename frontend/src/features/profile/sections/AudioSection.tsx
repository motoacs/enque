import { useTranslation } from "react-i18next";
import { useProfileStore } from "@/stores/profileStore";

export function AudioSection() {
  const { t } = useTranslation();
  const { editingProfile: p, updateEditingProfile: update } = useProfileStore();
  if (!p) return null;

  const isPreset = p.is_preset;

  return (
    <section>
      <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">
        {t("profile.audio")}
      </h3>
      <div className="space-y-2">
        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Audio Mode</label>
          <select
            value={p.audio_mode}
            onChange={(e) => update({ audio_mode: e.target.value })}
            disabled={isPreset}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
          >
            <option value="copy">Copy</option>
            <option value="aac">AAC</option>
            <option value="opus">Opus</option>
          </select>
        </div>

        {p.audio_mode !== "copy" && (
          <div className="flex items-center gap-2">
            <label className="text-xs text-zinc-400 w-28 shrink-0">Bitrate (kbps)</label>
            <input
              type="number"
              value={p.audio_bitrate}
              onChange={(e) => update({ audio_bitrate: Number(e.target.value) })}
              disabled={isPreset}
              min={32}
              max={1024}
              className="w-24 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
            />
          </div>
        )}
      </div>
    </section>
  );
}
