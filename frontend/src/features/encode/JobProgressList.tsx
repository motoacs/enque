import { useTranslation } from "react-i18next";
import { useEncodeStore } from "@/stores/encodeStore";
import { JobProgressItem } from "./JobProgressItem";

interface JobProgressListProps {
  onSelectJob: (jobId: string) => void;
  selectedJobId: string;
}

export function JobProgressList({ onSelectJob, selectedJobId }: JobProgressListProps) {
  const { t } = useTranslation();
  const jobProgress = useEncodeStore((s) => s.jobProgress);
  const jobs = Object.values(jobProgress);

  return (
    <div className="flex-1 overflow-y-auto">
      {jobs.map((job) => (
        <JobProgressItem
          key={job.jobId}
          job={job}
          selected={job.jobId === selectedJobId}
          onClick={() => onSelectJob(job.jobId)}
        />
      ))}
      {jobs.length === 0 && (
        <div className="text-xs text-zinc-500 text-center py-8">
          {t("encode.pending")}
        </div>
      )}
    </div>
  );
}
