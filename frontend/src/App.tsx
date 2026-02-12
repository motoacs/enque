import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { TopBar } from "@/components/layout/TopBar";
import { QueuePanel } from "@/features/queue/QueuePanel";
import { ProfileEditor } from "@/features/profile/ProfileEditor";
import { OutputSettingsPanel } from "@/features/output/OutputSettingsPanel";
import { CommandPreview } from "@/features/preview/CommandPreview";
import { EncodePanel } from "@/features/encode/EncodePanel";
import { SettingsDialog } from "@/features/settings/SettingsDialog";
import { TempCleanupDialog } from "@/features/settings/TempCleanupDialog";
import { useEncodeStore } from "@/stores/encodeStore";
import { useEditStore } from "@/stores/editStore";
import { useProfileStore } from "@/stores/profileStore";
import { useAppStore } from "@/stores/appStore";
import { Play } from "lucide-react";
import * as api from "@/lib/api";
import { registerEventListeners } from "@/lib/events";
import i18n from "@/lib/i18n";

function App() {
  const { t } = useTranslation();
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [tempFiles, setTempFiles] = useState<string[]>([]);
  const [tempDialogOpen, setTempDialogOpen] = useState(false);
  const [encodeError, setEncodeError] = useState<string | null>(null);
  const sessionState = useEncodeStore((s) => s.sessionState);
  const jobs = useEditStore((s) => s.jobs);
  const outputSettings = useEditStore((s) => s.outputSettings);
  const editingProfile = useProfileStore((s) => s.editingProfile);
  const config = useAppStore((s) => s.config);
  const setConfig = useAppStore((s) => s.setConfig);
  const setTools = useAppStore((s) => s.setTools);
  const setProfiles = useProfileStore((s) => s.setProfiles);
  const setEditingProfile = useProfileStore((s) => s.setEditingProfile);
  const setSelectedProfileId = useEditStore((s) => s.setSelectedProfileId);
  const setOutputSettings = useEditStore((s) => s.setOutputSettings);

  useEffect(() => {
    registerEventListeners();

    api.bootstrap().then((result: any) => {
      if (!result) return;

      if (result.config) {
        setConfig(result.config);
        i18n.changeLanguage(result.config.language || "ja");

        setOutputSettings({
          outputFolderMode: result.config.output_folder_mode,
          outputFolderPath: result.config.output_folder_path,
          outputNameTemplate: result.config.output_name_template,
          overwriteMode: result.config.overwrite_mode,
        });
      }
      if (result.profiles) {
        setProfiles(result.profiles);
        if (result.config?.default_profile_id) {
          const defaultProfile = result.profiles.find(
            (p: any) => p.id === result.config.default_profile_id
          );
          if (defaultProfile) {
            setEditingProfile(defaultProfile);
            setSelectedProfileId(defaultProfile.id);
          }
        } else if (result.profiles.length > 0) {
          setEditingProfile(result.profiles[0]);
          setSelectedProfileId(result.profiles[0].id);
        }
      }
      if (result.tools) {
        setTools(result.tools);
      }

      if (result.temp_artifacts && result.temp_artifacts.length > 0) {
        setTempFiles(result.temp_artifacts);
        setTempDialogOpen(true);
      }
    }).catch((err: unknown) => {
      console.error("Bootstrap failed:", err);
    });
  }, []);

  const isEncoding = sessionState === "running" || sessionState === "stopping" || sessionState === "aborting";
  const canStart = jobs.length > 0 && sessionState === "idle" && editingProfile !== null;

  const handleStartEncode = async () => {
    if (!editingProfile || !config) return;

    const request = {
      jobs: jobs.map((j) => ({
        job_id: j.jobId,
        input_path: j.inputPath,
      })),
      profile: editingProfile,
      app_config_snapshot: {
        max_concurrent_jobs: config.max_concurrent_jobs,
        on_error: config.on_error,
        decoder_fallback: config.decoder_fallback,
        keep_failed_temp: config.keep_failed_temp,
        no_output_timeout_sec: config.no_output_timeout_sec,
        no_progress_timeout_sec: config.no_progress_timeout_sec,
        post_complete_action: config.post_complete_action,
        post_complete_command: config.post_complete_command,
        output_folder_mode: outputSettings.outputFolderMode,
        output_folder_path: outputSettings.outputFolderPath,
        output_name_template: outputSettings.outputNameTemplate,
        overwrite_mode: outputSettings.overwriteMode,
        nvencc_path: config.nvencc_path,
      },
    };

    try {
      setEncodeError(null);
      await api.startEncode(request);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      console.error("Failed to start encode:", msg);
      setEncodeError(msg);
    }
  };

  return (
    <div className="flex flex-col h-screen font-body" style={{ background: '#0a0a0f' }}>
      <TopBar onSettingsClick={() => setSettingsOpen(true)} />

      <main className="flex-1 flex overflow-hidden">
        {isEncoding || sessionState === "completed" || sessionState === "aborted" ? (
          <div className="flex-1 flex flex-col overflow-hidden">
            <EncodePanel />
          </div>
        ) : (
          <>
            {/* Left panel: Queue */}
            <div className="w-72 flex flex-col" style={{ borderRight: '1px solid rgba(255,255,255,0.06)', background: 'rgba(18, 18, 26, 0.5)' }}>
              <QueuePanel />
            </div>

            {/* Center panel: Profile Editor + Output + Preview */}
            <div className="flex-1 flex flex-col overflow-hidden">
              <div className="flex-1 overflow-y-auto min-h-0">
                <ProfileEditor />
              </div>
              <OutputSettingsPanel />
              <CommandPreview command="" />
            </div>
          </>
        )}
      </main>

      {/* Bottom bar: Encode start */}
      {!isEncoding && sessionState === "idle" && (
        <footer className="shrink-0" style={{ borderTop: '1px solid rgba(255,255,255,0.06)', background: 'rgba(18, 18, 26, 0.5)' }}>
          {encodeError && (
            <div
              className="px-5 py-2 text-xs"
              style={{ color: '#fbbf24', background: 'rgba(251, 191, 36, 0.08)', borderBottom: '1px solid rgba(251, 191, 36, 0.15)' }}
            >
              {t("encode.startFailed", { error: encodeError })}
            </div>
          )}
          <div className="px-5 py-3 flex items-center justify-between">
            <div className="text-xs font-mono" style={{ color: '#5c5c68' }}>
              {jobs.length > 0
                ? t("queue.items", { count: jobs.length })
                : ""}
            </div>
            <button
              onClick={handleStartEncode}
              disabled={!canStart}
              className={`btn-primary ${canStart ? 'animate-pulse-glow' : ''}`}
            >
              <Play size={14} />
              {t("encode.start")}
            </button>
          </div>
        </footer>
      )}

      <SettingsDialog
        open={settingsOpen}
        onClose={() => setSettingsOpen(false)}
      />

      <TempCleanupDialog
        open={tempDialogOpen}
        tempFiles={tempFiles}
        onCleanup={() => {
          api.cleanupTempArtifacts(tempFiles).catch(console.error);
          setTempDialogOpen(false);
          setTempFiles([]);
        }}
        onDismiss={() => {
          setTempDialogOpen(false);
        }}
      />
    </div>
  );
}

export default App;
