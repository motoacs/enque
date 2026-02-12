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
  const sessionState = useEncodeStore((s) => s.sessionState);
  const jobs = useEditStore((s) => s.jobs);
  const outputSettings = useEditStore((s) => s.outputSettings);
  const editingProfile = useProfileStore((s) => s.editingProfile);
  const config = useAppStore((s) => s.config);
  const setConfig = useAppStore((s) => s.setConfig);
  const setTools = useAppStore((s) => s.setTools);
  const setProfiles = useProfileStore((s) => s.setProfiles);
  const setEditingProfile = useProfileStore((s) => s.setEditingProfile);
  const setOutputSettings = useEditStore((s) => s.setOutputSettings);

  useEffect(() => {
    registerEventListeners();

    api.bootstrap().then((result: any) => {
      if (!result) return;

      // Set stores from bootstrap data
      if (result.config) {
        setConfig(result.config);
        i18n.changeLanguage(result.config.language || "ja");

        // Reflect output defaults from config
        setOutputSettings({
          outputFolderMode: result.config.output_folder_mode,
          outputFolderPath: result.config.output_folder_path,
          outputNameTemplate: result.config.output_name_template,
          outputContainer: result.config.output_container,
          overwriteMode: result.config.overwrite_mode,
        });
      }
      if (result.profiles) {
        setProfiles(result.profiles);
        // Select default profile if configured
        if (result.config?.default_profile_id) {
          const defaultProfile = result.profiles.find(
            (p: any) => p.id === result.config.default_profile_id
          );
          if (defaultProfile) {
            setEditingProfile(defaultProfile);
          }
        } else if (result.profiles.length > 0) {
          setEditingProfile(result.profiles[0]);
        }
      }
      if (result.tools) {
        setTools(result.tools);
      }

      // Show temp cleanup dialog if there are leftover temp files
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
        output_container: outputSettings.outputContainer,
        overwrite_mode: outputSettings.overwriteMode,
        nvencc_path: config.nvencc_path,
      },
    };

    try {
      await api.startEncode(request);
    } catch (err) {
      console.error("Failed to start encode:", err);
    }
  };

  return (
    <div className="flex flex-col h-screen bg-zinc-900 text-zinc-100">
      <TopBar onSettingsClick={() => setSettingsOpen(true)} />

      <main className="flex-1 flex overflow-hidden">
        {isEncoding || sessionState === "completed" || sessionState === "aborted" ? (
          /* Encode monitoring view */
          <div className="flex-1 flex flex-col overflow-hidden">
            <EncodePanel />
          </div>
        ) : (
          /* Normal editing view */
          <>
            {/* Left panel: Queue */}
            <div className="w-72 border-r border-zinc-700 flex flex-col">
              <QueuePanel />
            </div>

            {/* Center panel: Profile Editor */}
            <div className="flex-1 flex flex-col overflow-hidden">
              <div className="flex-1 overflow-y-auto">
                <ProfileEditor />
              </div>
              <OutputSettingsPanel />
              <CommandPreview command="" />
            </div>
          </>
        )}
      </main>

      {/* Bottom bar: Encode controls */}
      {!isEncoding && sessionState === "idle" && (
        <footer className="px-4 py-2 border-t border-zinc-700 bg-zinc-800 flex items-center justify-between shrink-0">
          <div className="text-xs text-zinc-500">
            {jobs.length > 0
              ? t("queue.items", { count: jobs.length })
              : ""}
          </div>
          <button
            onClick={handleStartEncode}
            disabled={!canStart}
            className="flex items-center gap-2 px-4 py-1.5 rounded bg-blue-600 hover:bg-blue-500 disabled:bg-zinc-700 disabled:text-zinc-500 text-white text-sm font-medium transition-colors"
          >
            <Play size={14} />
            {t("encode.start")}
          </button>
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
