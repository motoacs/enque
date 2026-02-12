import { useTranslation } from "react-i18next";
import { useProfileStore } from "@/stores/profileStore";

export function SpeedSection() {
  const { t } = useTranslation();
  const { editingProfile: p, updateEditingProfile: update } = useProfileStore();
  if (!p) return null;

  const isPreset = p.is_preset;

  return (
    <section>
      <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">
        {t("profile.speed")}
      </h3>
      <div className="space-y-2">
        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Split Enc</label>
          <select
            value={p.split_enc}
            onChange={(e) => update({ split_enc: e.target.value })}
            disabled={isPreset}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
          >
            {["off", "auto", "auto_forced", "forced_2", "forced_3", "forced_4"].map((v) => (
              <option key={v} value={v}>{v}</option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Parallel</label>
          <select
            value={p.parallel}
            onChange={(e) => update({ parallel: e.target.value })}
            disabled={isPreset}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
          >
            {["off", "auto", "2", "3"].map((v) => (
              <option key={v} value={v}>{v}</option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Decoder</label>
          <select
            value={p.decoder}
            onChange={(e) => update({ decoder: e.target.value })}
            disabled={isPreset}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
          >
            <option value="avhw">avhw (Hardware)</option>
            <option value="avsw">avsw (Software)</option>
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Device</label>
          <select
            value={p.device}
            onChange={(e) => update({ device: e.target.value })}
            disabled={isPreset}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
          >
            <option value="auto">Auto</option>
            {[0, 1, 2, 3].map((id) => (
              <option key={id} value={String(id)}>GPU {id}</option>
            ))}
          </select>
        </div>
      </div>
    </section>
  );
}
