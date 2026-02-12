import { useTranslation } from "react-i18next";
import type { JobProgress } from "@/stores/encodeStore";
import { CheckCircle, XCircle, Clock, Loader, MinusCircle, AlertTriangle } from "lucide-react";

interface JobProgressItemProps {
  job: JobProgress;
  selected: boolean;
  onClick: () => void;
}

export function JobProgressItem({ job, selected, onClick }: JobProgressItemProps) {
  const { t } = useTranslation();

  const fileName = job.inputPath ? job.inputPath.split(/[\\/]/).pop() || job.jobId : job.jobId;
  const percent = job.percent ?? 0;

  const statusIcon = () => {
    switch (job.status) {
      case "completed":
        return <CheckCircle size={14} className="text-green-400" />;
      case "failed":
        return <XCircle size={14} className="text-red-400" />;
      case "cancelled":
        return <MinusCircle size={14} className="text-yellow-400" />;
      case "timeout":
        return <AlertTriangle size={14} className="text-orange-400" />;
      case "skipped":
        return <MinusCircle size={14} className="text-zinc-500" />;
      case "running":
        return <Loader size={14} className="text-blue-400 animate-spin" />;
      default:
        return <Clock size={14} className="text-zinc-500" />;
    }
  };

  const formatETA = (sec?: number) => {
    if (sec == null) return "";
    const h = Math.floor(sec / 3600);
    const m = Math.floor((sec % 3600) / 60);
    const s = Math.floor(sec % 60);
    if (h > 0) return `${h}:${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
    return `${m}:${String(s).padStart(2, "0")}`;
  };

  return (
    <div
      onClick={onClick}
      className={`px-3 py-2 cursor-pointer border-b border-zinc-700/50 hover:bg-zinc-700/30 ${
        selected ? "bg-zinc-700/50" : ""
      }`}
    >
      <div className="flex items-center gap-2 mb-1">
        {statusIcon()}
        <span className="text-xs text-zinc-200 truncate flex-1" title={job.inputPath}>
          {fileName}
        </span>
        {job.status === "running" && job.fps != null && (
          <span className="text-[10px] text-zinc-500">
            {job.fps.toFixed(1)} fps
          </span>
        )}
      </div>

      {job.status === "running" && (
        <>
          <div className="h-1.5 bg-zinc-700 rounded overflow-hidden mb-1">
            <div
              className="h-full bg-blue-500 transition-all duration-300"
              style={{ width: `${percent}%` }}
            />
          </div>
          <div className="flex justify-between text-[10px] text-zinc-500">
            <span>{percent.toFixed(1)}%</span>
            <div className="flex gap-2">
              {job.bitrateKbps != null && (
                <span>{job.bitrateKbps.toFixed(0)} kb/s</span>
              )}
              {job.etaSec != null && (
                <span>{t("encode.eta")} {formatETA(job.etaSec)}</span>
              )}
            </div>
          </div>
        </>
      )}

      {job.status === "failed" && job.errorMessage && (
        <div className="text-[10px] text-red-400 truncate mt-0.5" title={job.errorMessage}>
          {job.errorMessage}
        </div>
      )}
    </div>
  );
}
