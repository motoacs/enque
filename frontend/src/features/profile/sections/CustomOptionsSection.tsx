import { useTranslation } from "react-i18next";
import { useProfileStore } from "@/stores/profileStore";

export function CustomOptionsSection() {
  const { t } = useTranslation();
  const { editingProfile: p, updateEditingProfile: update } = useProfileStore();
  if (!p) return null;

  return (
    <section>
      <h3 className="section-heading mb-3">
        {t("profile.customOptions")}
      </h3>
      <textarea
        value={p.custom_options}
        onChange={(e) => update({ custom_options: e.target.value })}
        disabled={p.is_preset}
        placeholder="--vpp-nlmeans sigma=0.005 --gop-len 300"
        rows={3}
        className="w-full form-input font-mono resize-y leading-relaxed"
      />
      <p className="text-[10px] mt-1.5 font-mono" style={{ color: '#3c3c48' }}>
        {p.custom_options.length} / 4096
      </p>
    </section>
  );
}
