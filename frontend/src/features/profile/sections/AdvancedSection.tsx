import { useTranslation } from "react-i18next";
import { useProfileStore } from "@/stores/profileStore";
import { useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";

export function AdvancedSection() {
  const { t } = useTranslation();
  const { editingProfile: p, updateAdvanced } = useProfileStore();
  const [expanded, setExpanded] = useState(false);

  if (!p) return null;

  const isPreset = p.is_preset;
  const adv = p.nvencc_advanced;

  const setStr = (field: string, value: string) => {
    updateAdvanced({ [field]: value } as any);
  };

  const setNullInt = (field: string, value: string) => {
    updateAdvanced({ [field]: value === "" ? null : Number(value) } as any);
  };

  const setBool = (field: string, value: boolean) => {
    updateAdvanced({ [field]: value } as any);
  };

  return (
    <section>
      <button
        onClick={() => setExpanded(!expanded)}
        className="flex items-center gap-1 text-xs font-semibold text-zinc-400 uppercase mb-2 hover:text-zinc-300"
      >
        {expanded ? <ChevronDown size={12} /> : <ChevronRight size={12} />}
        {t("profile.advanced")}
      </button>

      {expanded && (
        <div className="space-y-2 ml-1">
          {[
            { label: "Interlace", field: "interlace", type: "text" },
            { label: "SW Decoder", field: "avsw_decoder", type: "text" },
            { label: "Input CSP", field: "input_csp", type: "text" },
            { label: "Output CSP", field: "output_csp", type: "text" },
            { label: "Tune", field: "tune", type: "text" },
            { label: "MV Precision", field: "mv_precision", type: "text" },
            { label: "Level", field: "level", type: "text" },
            { label: "Profile", field: "profile", type: "text" },
            { label: "Tier", field: "tier", type: "text" },
            { label: "Trim", field: "trim", type: "text" },
            { label: "Seek", field: "seek", type: "text" },
            { label: "Seek To", field: "seekto", type: "text" },
          ].map(({ label, field }) => (
            <div key={field} className="flex items-center gap-2">
              <label className="text-xs text-zinc-400 w-28 shrink-0">{label}</label>
              <input
                type="text"
                value={(adv as any)[field] || ""}
                onChange={(e) => setStr(field, e.target.value)}
                disabled={isPreset}
                className="flex-1 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
              />
            </div>
          ))}

          {[
            { label: "Max Bitrate", field: "max_bitrate" },
            { label: "VBR Quality", field: "vbr_quality" },
            { label: "Lookahead Lv", field: "lookahead_level" },
            { label: "Refs Forward", field: "refs_forward" },
            { label: "Refs Backward", field: "refs_backward" },
            { label: "Output Thread", field: "output_thread" },
          ].map(({ label, field }) => (
            <div key={field} className="flex items-center gap-2">
              <label className="text-xs text-zinc-400 w-28 shrink-0">{label}</label>
              <input
                type="number"
                value={(adv as any)[field] ?? ""}
                onChange={(e) => setNullInt(field, e.target.value)}
                disabled={isPreset}
                placeholder="auto"
                className="w-24 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60 placeholder:text-zinc-600"
              />
            </div>
          ))}

          <div className="flex items-center gap-4">
            {[
              { label: "WeightP", field: "weightp" },
              { label: "SSIM", field: "ssim" },
              { label: "PSNR", field: "psnr" },
            ].map(({ label, field }) => (
              <label key={field} className="flex items-center gap-1.5 text-xs text-zinc-400">
                <input
                  type="checkbox"
                  checked={(adv as any)[field] || false}
                  onChange={(e) => setBool(field, e.target.checked)}
                  disabled={isPreset}
                  className="rounded"
                />
                {label}
              </label>
            ))}
          </div>

          {/* Advanced raw audio/sub/data overrides */}
          <div className="mt-2 pt-2 border-t border-zinc-700">
            <p className="text-xs text-zinc-500 mb-2">{t("profile.advancedOverrides")}</p>
            {[
              { label: "Video Meta", field: "video_metadata" },
              { label: "Audio Copy", field: "audio_copy" },
              { label: "Audio Codec", field: "audio_codec" },
              { label: "Audio BR", field: "audio_bitrate" },
              { label: "Audio Quality", field: "audio_quality" },
              { label: "Audio SR", field: "audio_samplerate" },
              { label: "Audio Meta", field: "audio_metadata" },
              { label: "Sub Copy", field: "sub_copy" },
              { label: "Sub Meta", field: "sub_metadata" },
              { label: "Data Copy", field: "data_copy" },
              { label: "Attach Copy", field: "attachment_copy" },
              { label: "Metadata", field: "metadata" },
            ].map(({ label, field }) => (
              <div key={field} className="flex items-center gap-2 mb-1">
                <label className="text-xs text-zinc-500 w-28 shrink-0">{label}</label>
                <input
                  type="text"
                  value={(adv as any)[field] || ""}
                  onChange={(e) => setStr(field, e.target.value)}
                  disabled={isPreset}
                  className="flex-1 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
                />
              </div>
            ))}
          </div>
        </div>
      )}
    </section>
  );
}
