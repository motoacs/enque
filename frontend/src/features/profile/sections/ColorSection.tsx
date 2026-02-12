import { useTranslation } from "react-i18next";
import { useProfileStore } from "@/stores/profileStore";

const colormatrixOpts = ["auto", "bt709", "bt2020nc", "bt2020c", "smpte240m", "fcc", "GBR"];
const transferOpts = ["auto", "bt709", "smpte2084", "arib-std-b67", "smpte240m", "linear", "log100"];
const colorprimOpts = ["auto", "bt709", "bt2020", "smpte240m", "smpte431", "smpte432"];
const colorrangeOpts = ["auto", "limited", "full"];
const dhdrOpts = ["off", "copy"];

export function ColorSection() {
  const { t } = useTranslation();
  const { editingProfile: p, updateEditingProfile: update } = useProfileStore();
  if (!p) return null;

  const isPreset = p.is_preset;

  const selects = [
    { label: "Colormatrix", field: "colormatrix" as const, opts: colormatrixOpts },
    { label: "Transfer", field: "transfer" as const, opts: transferOpts },
    { label: "Color Prim", field: "colorprim" as const, opts: colorprimOpts },
    { label: "Color Range", field: "colorrange" as const, opts: colorrangeOpts },
    { label: "HDR10+", field: "dhdr10_info" as const, opts: dhdrOpts },
  ];

  return (
    <section>
      <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">
        {t("profile.color")}
      </h3>
      <div className="space-y-2">
        {selects.map(({ label, field, opts }) => (
          <div key={field} className="flex items-center gap-2">
            <label className="text-xs text-zinc-400 w-28 shrink-0">{label}</label>
            <select
              value={p[field]}
              onChange={(e) => update({ [field]: e.target.value })}
              disabled={isPreset}
              className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 disabled:opacity-60"
            >
              {opts.map((v) => (
                <option key={v} value={v}>{v}</option>
              ))}
            </select>
          </div>
        ))}
      </div>
    </section>
  );
}
