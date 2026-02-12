import { useTranslation } from "react-i18next";
import { useProfileStore } from "@/stores/profileStore";

export function VideoBasicSection() {
  const { t } = useTranslation();
  const { editingProfile: p, updateEditingProfile: update } = useProfileStore();
  if (!p) return null;

  const isPreset = p.is_preset;

  return (
    <section>
      <h3 className="section-heading mb-3">
        {t("profile.video")}
      </h3>
      <div className="space-y-2.5">
        <div className="flex items-center gap-2">
          <label className="form-label w-28">Codec</label>
          <div className="flex gap-1">
            {["h264", "hevc", "av1"].map((c) => (
              <button
                key={c}
                onClick={() => !isPreset && update({ codec: c })}
                className={`codec-btn ${
                  p.codec === c ? "codec-btn-active" : "codec-btn-inactive"
                } ${isPreset ? "opacity-45 cursor-not-allowed" : ""}`}
              >
                {c.toUpperCase()}
              </button>
            ))}
          </div>
        </div>

        <div className="flex items-center gap-2">
          <label className="form-label w-28">Rate Control</label>
          <select
            value={p.rate_control}
            onChange={(e) => update({ rate_control: e.target.value })}
            disabled={isPreset}
            className="form-input"
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
            className="w-20 form-input font-mono"
          />
        </div>

        <div className="flex items-center gap-2">
          <label className="form-label w-28">Preset</label>
          <select
            value={p.preset}
            onChange={(e) => update({ preset: e.target.value })}
            disabled={isPreset}
            className="form-input"
          >
            {["P1", "P2", "P3", "P4", "P5", "P6", "P7"].map((pr) => (
              <option key={pr} value={pr}>{pr}</option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="form-label w-28">Output Depth</label>
          <select
            value={p.output_depth}
            onChange={(e) => update({ output_depth: Number(e.target.value) })}
            disabled={isPreset}
            className="form-input"
          >
            <option value={8}>8-bit</option>
            <option value={10}>10-bit</option>
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="form-label w-28">Multipass</label>
          <select
            value={p.multipass}
            onChange={(e) => update({ multipass: e.target.value })}
            disabled={isPreset}
            className="form-input"
          >
            {["none", "quarter", "full"].map((mp) => (
              <option key={mp} value={mp}>{mp}</option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="form-label w-28">Output Res</label>
          <input
            type="text"
            value={p.output_res}
            onChange={(e) => update({ output_res: e.target.value })}
            disabled={isPreset}
            placeholder="e.g. 1920x1080"
            className="flex-1 form-input font-mono"
          />
        </div>

        <div className="flex items-center gap-2">
          <label className="form-label w-28">{t("output.container")}</label>
          <select
            value={p.output_container}
            onChange={(e) => update({ output_container: e.target.value })}
            disabled={isPreset}
            className="form-input"
          >
            {["mp4", "mkv", "mov", "webm"].map((c) => (
              <option key={c} value={c}>{c}</option>
            ))}
          </select>
        </div>
      </div>
    </section>
  );
}
