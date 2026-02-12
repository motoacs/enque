import { useTranslation } from "react-i18next";
import { useAppStore, type AppConfig } from "@/stores/appStore";
import { X } from "lucide-react";

interface SettingsDialogProps {
  open: boolean;
  onClose: () => void;
}

export function SettingsDialog({ open, onClose }: SettingsDialogProps) {
  const { t, i18n } = useTranslation();
  const { config, updateConfig } = useAppStore();

  if (!open || !config) return null;

  const handleLanguageChange = (lang: string) => {
    updateConfig({ language: lang });
    i18n.changeLanguage(lang);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-zinc-800 rounded-lg shadow-xl w-[500px] max-h-[80vh] flex flex-col border border-zinc-700">
        <div className="flex items-center justify-between px-4 py-3 border-b border-zinc-700">
          <h2 className="text-sm font-semibold text-zinc-200">{t("settings.title")}</h2>
          <button onClick={onClose} className="text-zinc-400 hover:text-zinc-200">
            <X size={16} />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {/* Tool Paths */}
          <section>
            <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">{t("settings.tools")}</h3>
            {([
              ["nvencc_path", t("settings.nvenccPath")] as const,
              ["qsvenc_path", t("settings.qsvencPath")] as const,
              ["ffmpeg_path", t("settings.ffmpegPath")] as const,
              ["ffprobe_path", t("settings.ffprobePath")] as const,
            ] as const).map(([field, label]) => (
              <div key={field} className="flex items-center gap-2 mb-2">
                <label className="text-xs text-zinc-400 w-28 shrink-0">{label}</label>
                <input
                  type="text"
                  value={config[field]}
                  onChange={(e) => updateConfig({ [field]: e.target.value })}
                  className="flex-1 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600"
                />
              </div>
            ))}
          </section>

          {/* Execution */}
          <section>
            <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">{t("settings.execution")}</h3>
            <div className="space-y-2">
              <div className="flex items-center gap-2">
                <label className="text-xs text-zinc-400 w-28 shrink-0">{t("settings.maxJobs")}</label>
                <input
                  type="number"
                  value={config.max_concurrent_jobs}
                  onChange={(e) => updateConfig({ max_concurrent_jobs: Number(e.target.value) })}
                  min={1}
                  max={8}
                  className="w-16 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600"
                />
              </div>

              <div className="flex items-center gap-2">
                <label className="text-xs text-zinc-400 w-28 shrink-0">{t("settings.onError")}</label>
                <select
                  value={config.on_error}
                  onChange={(e) => updateConfig({ on_error: e.target.value })}
                  className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600"
                >
                  <option value="skip">{t("settings.skip")}</option>
                  <option value="stop">{t("settings.stop")}</option>
                </select>
              </div>

              <label className="flex items-center gap-2 text-xs text-zinc-400">
                <input
                  type="checkbox"
                  checked={config.decoder_fallback}
                  onChange={(e) => updateConfig({ decoder_fallback: e.target.checked })}
                />
                {t("settings.decoderFallback")}
              </label>

              <label className="flex items-center gap-2 text-xs text-zinc-400">
                <input
                  type="checkbox"
                  checked={config.keep_failed_temp}
                  onChange={(e) => updateConfig({ keep_failed_temp: e.target.checked })}
                />
                {t("settings.keepFailedTemp")}
              </label>

              <div className="flex items-center gap-2">
                <label className="text-xs text-zinc-400 w-28 shrink-0">{t("settings.noOutputTimeout")}</label>
                <input
                  type="number"
                  value={config.no_output_timeout_sec}
                  onChange={(e) => updateConfig({ no_output_timeout_sec: Number(e.target.value) })}
                  min={30}
                  max={86400}
                  className="w-20 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600"
                />
              </div>

              <div className="flex items-center gap-2">
                <label className="text-xs text-zinc-400 w-28 shrink-0">{t("settings.noProgressTimeout")}</label>
                <input
                  type="number"
                  value={config.no_progress_timeout_sec}
                  onChange={(e) => updateConfig({ no_progress_timeout_sec: Number(e.target.value) })}
                  min={30}
                  max={86400}
                  className="w-20 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600"
                />
              </div>
            </div>
          </section>

          {/* Post Action */}
          <section>
            <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">{t("settings.postAction")}</h3>
            <select
              value={config.post_complete_action}
              onChange={(e) => updateConfig({ post_complete_action: e.target.value })}
              className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 mb-2"
            >
              <option value="none">{t("settings.none")}</option>
              <option value="shutdown">{t("settings.shutdown")}</option>
              <option value="sleep">{t("settings.sleep")}</option>
              <option value="custom">{t("settings.custom")}</option>
            </select>
            {config.post_complete_action === "custom" && (
              <input
                type="text"
                value={config.post_complete_command}
                onChange={(e) => updateConfig({ post_complete_command: e.target.value })}
                placeholder="command..."
                className="w-full bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 font-mono"
              />
            )}
          </section>

          {/* Language */}
          <section>
            <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">{t("settings.language")}</h3>
            <select
              value={config.language}
              onChange={(e) => handleLanguageChange(e.target.value)}
              className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600"
            >
              <option value="ja">日本語</option>
              <option value="en">English</option>
            </select>
          </section>
        </div>

        <div className="px-4 py-3 border-t border-zinc-700 flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-1.5 text-xs bg-zinc-600 hover:bg-zinc-500 text-zinc-200 rounded transition-colors"
          >
            {t("common.close")}
          </button>
        </div>
      </div>
    </div>
  );
}
