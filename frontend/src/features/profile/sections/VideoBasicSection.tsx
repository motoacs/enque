import { useTranslation } from "react-i18next";
import { useProfileStore } from "@/stores/profileStore";

export function VideoBasicSection() {
  const { t } = useTranslation();
  const { editingProfile: p, updateEditingProfile: update } = useProfileStore();
  if (!p) return null;

  const isPreset = p.is_preset;

  return (
    <section>
      <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">
        {t("profile.video")}
      </h3>
      <div className="space-y-2">
        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Codec</label>
          <div className="flex gap-1">
            {["h264", "hevc", "av1"].map((c) => (
              <button
                key={c}
                onClick={() => !isPreset && update({ codec: c })}
                className={`px-3 py-1 text-xs rounded ${
                  p.codec === c
                    ? "bg-blue-600 text-white"
                    : "bg-zinc-700 text-zinc-300 hover:bg-zinc-600"
                } ${isPreset ? "opacity-60 cursor-not-allowed" : ""}`}
              >
                {c.toUpperCase()}
              </button>
            ))}
          </div>
        </div>

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Rate Control</label>
          <select
            value={p.rate_control}
            onChange={(e) => update({ rate_control: e.target.value })}
            disabled={isPreset}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
          >
            {["qvbr", "cqp", "cbr", "vbr"].map((rc) => (
              <option key={rc} value={rc}>{rc.toUpperCase()}</option>
            ))}
          </select>
          <input
            type="number"
            value={p.rate_value}
            onChange={(e) => update({ rate_value: Number(e.target.value) })}
            disabled={isPreset}
            className="w-20 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
          />
        </div>

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Preset</label>
          <select
            value={p.preset}
            onChange={(e) => update({ preset: e.target.value })}
            disabled={isPreset}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
          >
            {["P1", "P2", "P3", "P4", "P5", "P6", "P7"].map((pr) => (
              <option key={pr} value={pr}>{pr}</option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Output Depth</label>
          <select
            value={p.output_depth}
            onChange={(e) => update({ output_depth: Number(e.target.value) })}
            disabled={isPreset}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
          >
            <option value={8}>8-bit</option>
            <option value={10}>10-bit</option>
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Multipass</label>
          <select
            value={p.multipass}
            onChange={(e) => update({ multipass: e.target.value })}
            disabled={isPreset}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
          >
            {["none", "quarter", "full"].map((mp) => (
              <option key={mp} value={mp}>{mp}</option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-28 shrink-0">Output Res</label>
          <input
            type="text"
            value={p.output_res}
            onChange={(e) => update({ output_res: e.target.value })}
            disabled={isPreset}
            placeholder="e.g. 1920x1080"
            className="flex-1 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60 placeholder:text-zinc-600"
          />
        </div>
      </div>
    </section>
  );
}
