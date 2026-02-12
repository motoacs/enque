import { useTranslation } from "react-i18next";
import { useEditStore } from "@/stores/editStore";

const containers = ["mkv", "mp4", "mov", "webm"];

export function OutputSettingsPanel() {
  const { t } = useTranslation();
  const { outputSettings, setOutputSettings } = useEditStore();

  return (
    <section className="px-3 py-2 border-b border-zinc-700">
      <h3 className="text-xs font-semibold text-zinc-400 uppercase mb-2">
        {t("output.title")}
      </h3>
      <div className="space-y-2">
        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-24 shrink-0">{t("output.folder")}</label>
          <select
            value={outputSettings.outputFolderMode}
            onChange={(e) => setOutputSettings({ outputFolderMode: e.target.value })}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600"
          >
            <option value="same_as_input">{t("output.sameAsInput")}</option>
            <option value="specified">{t("output.specified")}</option>
          </select>
        </div>

        {outputSettings.outputFolderMode === "specified" && (
          <div className="flex items-center gap-2">
            <label className="text-xs text-zinc-400 w-24 shrink-0">Path</label>
            <input
              type="text"
              value={outputSettings.outputFolderPath}
              onChange={(e) => setOutputSettings({ outputFolderPath: e.target.value })}
              className="flex-1 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600"
            />
          </div>
        )}

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-24 shrink-0">{t("output.template")}</label>
          <input
            type="text"
            value={outputSettings.outputNameTemplate}
            onChange={(e) => setOutputSettings({ outputNameTemplate: e.target.value })}
            className="flex-1 bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600 font-mono"
          />
        </div>

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-24 shrink-0">{t("output.container")}</label>
          <select
            value={outputSettings.outputContainer}
            onChange={(e) => setOutputSettings({ outputContainer: e.target.value })}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600"
          >
            {containers.map((c) => (
              <option key={c} value={c}>{c}</option>
            ))}
          </select>
        </div>

        <div className="flex items-center gap-2">
          <label className="text-xs text-zinc-400 w-24 shrink-0">{t("output.overwrite")}</label>
          <select
            value={outputSettings.overwriteMode}
            onChange={(e) => setOutputSettings({ overwriteMode: e.target.value })}
            className="bg-zinc-700 text-zinc-200 text-xs rounded px-2 py-1 border border-zinc-600"
          >
            <option value="ask">{t("output.ask")}</option>
            <option value="auto_rename">{t("output.autoRename")}</option>
          </select>
        </div>
      </div>
    </section>
  );
}
