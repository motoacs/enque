import { useTranslation } from "react-i18next";
import { useEncodeStore } from "@/stores/encodeStore";
import { JobProgressItem } from "./JobProgressItem";
import * as api from "@/lib/api";

interface JobProgressListProps {
  onSelectJob: (jobId: string) => void;
  selectedJobId: string;
}

export function JobProgressList({ onSelectJob, selectedJobId }: JobProgressListProps) {
  const { t } = useTranslation();
  const jobProgress = useEncodeStore((s) => s.jobProgress);
  const sessionId = useEncodeStore((s) => s.sessionId);
  const skipPendingJob = useEncodeStore((s) => s.skipPendingJob);
  const jobs = Object.values(jobProgress);

  const handleSkip = async (jobId: string) => {
    try {
      await api.skipJob(sessionId, jobId);
      skipPendingJob(jobId);
    } catch (err) {
      console.error("Failed to skip job:", err);
    }
  };

  return (
    <div className="flex-1 overflow-y-auto">
      {jobs.map((job) => (
        <JobProgressItem
          key={job.jobId}
          job={job}
          selected={job.jobId === selectedJobId}
          onClick={() => onSelectJob(job.jobId)}
          onSkip={handleSkip}
        />
      ))}
      {jobs.length === 0 && (
        <div className="text-xs text-center py-8" style={{ color: '#3c3c48' }}>
          {t("encode.pending")}
        </div>
      )}
    </div>
  );
}
