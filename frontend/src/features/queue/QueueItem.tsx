import { X, Film } from "lucide-react";
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
    <div className="flex items-center gap-2.5 px-3 py-2 rounded-lg group transition-colors duration-150 hover:bg-white/[0.03]">
      <Film size={13} style={{ color: '#5c5c68', flexShrink: 0 }} />
      <div className="flex-1 min-w-0">
        <p className="text-xs truncate" style={{ color: '#e8e6e3' }}>{job.fileName}</p>
        {job.inputSizeBytes > 0 && (
          <p className="text-[10px] font-mono" style={{ color: '#5c5c68' }}>{formatSize(job.inputSizeBytes)}</p>
        )}
      </div>
      <button
        onClick={() => onRemove(job.jobId)}
        className="opacity-0 group-hover:opacity-100 p-1 rounded-md transition-all duration-150 hover:bg-white/[0.06]"
        style={{ color: '#5c5c68' }}
      >
        <X size={12} />
      </button>
    </div>
  );
}
