import { useTranslation } from "react-i18next";
import { Plus, Trash2, FolderOpen } from "lucide-react";
import { useEditStore, type QueueJob } from "@/stores/editStore";
import { useEncodeStore } from "@/stores/encodeStore";
import * as api from "@/lib/api";
import { DropZone } from "./DropZone";
import { QueueItem } from "./QueueItem";

let jobCounter = 0;

export function QueuePanel() {
  const { t } = useTranslation();
  const { jobs, addJobs, removeJob, clearJobs } = useEditStore();
  const sessionState = useEncodeStore((s) => s.sessionState);
  const isLocked = sessionState !== "idle";

  const handleFilesDropped = (paths: string[]) => {
    if (isLocked) return;
    const newJobs: QueueJob[] = paths.map((p) => ({
      jobId: `job-${++jobCounter}-${Date.now()}`,
      inputPath: p,
      inputSizeBytes: 0,
      fileName: p.split(/[/\\]/).pop() || p,
    }));
    addJobs(newJobs);
  };

  const handleAddFiles = async () => {
    if (isLocked) return;
    try {
      const paths = await api.openFileDialog();
      if (paths && paths.length > 0) {
        handleFilesDropped(paths);
      }
    } catch (err) {
      console.error("Failed to open file dialog:", err);
    }
  };

  const handleAddFolder = async () => {
    if (isLocked) return;
    try {
      const paths = await api.openDirectoryDialog();
      if (paths && paths.length > 0) {
        handleFilesDropped(paths);
      }
    } catch (err) {
      console.error("Failed to open directory dialog:", err);
    }
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between px-4 py-2.5" style={{ borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
        <h2 className="font-display text-xs font-semibold uppercase tracking-wider" style={{ color: '#9d9da7' }}>
          {t("queue.title")}
          {jobs.length > 0 && (
            <span className="ml-2 font-mono text-[10px]" style={{ color: '#e8a849' }}>
              {jobs.length}
            </span>
          )}
        </h2>
        {!isLocked && (
          <div className="flex items-center gap-0.5">
            <button
              onClick={handleAddFiles}
              className="icon-btn"
              title={t("queue.addFiles")}
            >
              <Plus size={14} />
            </button>
            <button
              onClick={handleAddFolder}
              className="icon-btn"
              title={t("queue.addFolder")}
            >
              <FolderOpen size={14} />
            </button>
            {jobs.length > 0 && (
              <button
                onClick={clearJobs}
                className="icon-btn hover:!text-red-400"
                title={t("queue.clear")}
              >
                <Trash2 size={14} />
              </button>
            )}
          </div>
        )}
      </div>

      <DropZone onFilesDropped={handleFilesDropped}>
        <div className="flex-1 overflow-y-auto p-1.5">
          {jobs.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full p-6 text-center gap-3">
              <div className="w-10 h-10 rounded-lg flex items-center justify-center" style={{ background: 'rgba(232, 168, 73, 0.08)', border: '1px dashed rgba(232, 168, 73, 0.2)' }}>
                <Plus size={18} style={{ color: '#e8a849', opacity: 0.6 }} />
              </div>
              <p className="text-xs" style={{ color: '#5c5c68' }}>
                {t("queue.empty")}
              </p>
            </div>
          ) : (
            <div className="space-y-0.5">
              {jobs.map((job) => (
                <QueueItem
                  key={job.jobId}
                  job={job}
                  onRemove={isLocked ? () => {} : removeJob}
                />
              ))}
            </div>
          )}
        </div>
      </DropZone>
    </div>
  );
}
