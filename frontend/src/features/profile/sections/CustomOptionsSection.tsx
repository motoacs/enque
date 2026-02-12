import { useTranslation } from "react-i18next";
import { useProfileStore } from "@/stores/profileStore";

export function CustomOptionsSection() {
  const { t } = useTranslation();
  const { editingProfile: p, updateEditingProfile: update } = useProfileStore();
  if (!p) return null;

  return (
    <section>
      <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">
        {t("profile.customOptions")}
      </h3>
      <textarea
        value={p.custom_options}
        onChange={(e) => update({ custom_options: e.target.value })}
        disabled={p.is_preset}
        placeholder="--vpp-nlmeans sigma=0.005 --gop-len 300"
        rows={3}
        className="w-full bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1.5 border border-zinc-600 disabled:opacity-60 placeholder:text-zinc-600 font-mono resize-y"
      />
      <p className="text-xs text-zinc-500 mt-1">
        {p.custom_options.length} / 4096
      </p>
    </section>
  );
}
