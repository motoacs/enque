import { X } from "lucide-react";
import type { QueueJob } from "@/stores/editStore";

interface QueueItemProps {
  job: QueueJob;
  onRemove: (jobId: string) => void;
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 * 1024 * 1024)
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`;
}

export function QueueItem({ job, onRemove }: QueueItemProps) {
  return (
    <div className="flex items-center gap-3 px-3 py-2 rounded hover:bg-zinc-700/50 group">
      <div className="flex-1 min-w-0">
        <p className="text-sm text-zinc-200 truncate">{job.fileName}</p>
        <p className="text-xs text-zinc-500">{formatSize(job.inputSizeBytes)}</p>
      </div>
      <button
        onClick={() => onRemove(job.jobId)}
        className="opacity-0 group-hover:opacity-100 p-1 rounded hover:bg-zinc-600 text-zinc-400 hover:text-zinc-200 transition-all"
      >
        <X size={14} />
      </button>
    </div>
  );
}
