import { useTranslation } from "react-i18next";
import { Plus, Trash2, FolderOpen } from "lucide-react";
import { useEditStore, type QueueJob } from "@/stores/editStore";
import { useEncodeStore } from "@/stores/encodeStore";
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

  const handleAddFiles = () => {
    // In production, this would call Wails file dialog
    // For now, this is a stub
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between px-3 py-2 border-b border-zinc-700">
        <h2 className="text-sm font-medium text-zinc-300">
          {t("queue.title")}
          {jobs.length > 0 && (
            <span className="ml-2 text-zinc-500">({jobs.length})</span>
          )}
        </h2>
        {!isLocked && (
          <div className="flex items-center gap-1">
            <button
              onClick={handleAddFiles}
              className="p-1.5 rounded hover:bg-zinc-700 text-zinc-400 hover:text-zinc-200"
              title={t("queue.addFiles")}
            >
              <Plus size={14} />
            </button>
            <button
              onClick={handleAddFiles}
              className="p-1.5 rounded hover:bg-zinc-700 text-zinc-400 hover:text-zinc-200"
              title={t("queue.addFolder")}
            >
              <FolderOpen size={14} />
            </button>
            {jobs.length > 0 && (
              <button
                onClick={clearJobs}
                className="p-1.5 rounded hover:bg-zinc-700 text-zinc-400 hover:text-red-400"
                title={t("queue.clear")}
              >
                <Trash2 size={14} />
              </button>
            )}
          </div>
        )}
      </div>

      <DropZone onFilesDropped={handleFilesDropped}>
        <div className="flex-1 overflow-y-auto p-1">
          {jobs.length === 0 ? (
            <div className="flex items-center justify-center h-full text-zinc-500 text-sm p-4 text-center">
              {t("queue.empty")}
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
