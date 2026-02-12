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
        return <CheckCircle size={13} style={{ color: '#34d399' }} />;
      case "failed":
        return <XCircle size={13} style={{ color: '#f87171' }} />;
      case "cancelled":
        return <MinusCircle size={13} style={{ color: '#fbbf24' }} />;
      case "timeout":
        return <AlertTriangle size={13} style={{ color: '#fb923c' }} />;
      case "skipped":
        return <MinusCircle size={13} style={{ color: '#5c5c68' }} />;
      case "running":
        return <Loader size={13} className="animate-spin" style={{ color: '#4ecdc4' }} />;
      default:
        return <Clock size={13} style={{ color: '#3c3c48' }} />;
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
      className="px-3.5 py-2.5 cursor-pointer transition-colors duration-100"
      style={{
        borderBottom: '1px solid rgba(255,255,255,0.04)',
        background: selected ? 'rgba(78, 205, 196, 0.06)' : 'transparent',
        borderLeft: selected ? '2px solid #4ecdc4' : '2px solid transparent',
      }}
    >
      <div className="flex items-center gap-2 mb-1">
        {statusIcon()}
        <span className="text-xs truncate flex-1" style={{ color: '#e8e6e3' }} title={job.inputPath}>
          {fileName}
        </span>
        {job.status === "running" && job.fps != null && (
          <span className="text-[10px] font-mono" style={{ color: '#4ecdc4' }}>
            {job.fps.toFixed(1)} fps
          </span>
        )}
      </div>

      {job.status === "running" && (
        <>
          <div className="progress-track mb-1" style={{ height: '4px' }}>
            <div
              className="progress-fill-sm"
              style={{ width: `${percent}%` }}
            />
          </div>
          <div className="flex justify-between text-[10px] font-mono" style={{ color: '#5c5c68' }}>
            <span>{percent.toFixed(1)}%</span>
            <div className="flex gap-3">
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
        <div className="text-[10px] truncate mt-0.5" style={{ color: '#f87171' }} title={job.errorMessage}>
          {job.errorMessage}
        </div>
      )}
    </div>
  );
}
