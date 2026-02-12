import { useTranslation } from "react-i18next";
import { useEditStore } from "@/stores/editStore";

export function OutputSettingsPanel() {
  const { t } = useTranslation();
  const { outputSettings, setOutputSettings } = useEditStore();

  return (
    <section className="shrink-0 px-4 py-3" style={{ borderTop: '1px solid rgba(255,255,255,0.08)', background: 'rgba(12, 12, 18, 0.7)' }}>
      <h3 className="section-heading mb-3">
        {t("output.title")}
      </h3>
      <div className="space-y-2.5">
        <div className="flex items-center gap-2">
          <label className="form-label w-24">{t("output.folder")}</label>
          <select
            value={outputSettings.outputFolderMode}
            onChange={(e) => setOutputSettings({ outputFolderMode: e.target.value })}
            className="form-input"
          >
            <option value="same_as_input">{t("output.sameAsInput")}</option>
            <option value="specified">{t("output.specified")}</option>
          </select>
        </div>

        {outputSettings.outputFolderMode === "specified" && (
          <div className="flex items-center gap-2">
            <label className="form-label w-24">Path</label>
            <input
              type="text"
              value={outputSettings.outputFolderPath}
              onChange={(e) => setOutputSettings({ outputFolderPath: e.target.value })}
              className="flex-1 form-input font-mono"
            />
          </div>
        )}

        <div className="flex items-center gap-2">
          <label className="form-label w-24">{t("output.template")}</label>
          <input
            type="text"
            value={outputSettings.outputNameTemplate}
            onChange={(e) => setOutputSettings({ outputNameTemplate: e.target.value })}
            className="flex-1 form-input font-mono"
          />
        </div>

        <div className="flex items-center gap-2">
          <label className="form-label w-24">{t("output.overwrite")}</label>
          <select
            value={outputSettings.overwriteMode}
            onChange={(e) => setOutputSettings({ overwriteMode: e.target.value })}
            className="form-input"
          >
            <option value="ask">{t("output.ask")}</option>
            <option value="auto_rename">{t("output.autoRename")}</option>
          </select>
        </div>
      </div>
    </section>
  );
}
