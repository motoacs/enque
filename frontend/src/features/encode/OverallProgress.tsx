import { useTranslation } from "react-i18next";
import { useEncodeStore } from "@/stores/encodeStore";

export function OverallProgress() {
  const { t } = useTranslation();
  const jobProgress = useEncodeStore((s) => s.jobProgress);

  const jobs = Object.values(jobProgress);
  const total = jobs.length;
  const completed = jobs.filter(
    (j) => j.status === "completed" || j.status === "failed" || j.status === "cancelled" || j.status === "timeout" || j.status === "skipped"
  ).length;

  const percent = total > 0 ? Math.round((completed / total) * 100) : 0;

  return (
    <div className="px-5 py-3 space-y-2" style={{ borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
      <div className="flex justify-between items-baseline">
        <span className="font-display text-xs font-semibold uppercase tracking-wider" style={{ color: '#4ecdc4' }}>
          {t("encode.progress")}
        </span>
        <span className="text-xs font-mono" style={{ color: '#9d9da7' }}>
          {completed} / {total}
          <span className="ml-2" style={{ color: '#4ecdc4' }}>{percent}%</span>
        </span>
      </div>
      <div className="progress-track" style={{ height: '8px' }}>
        <div
          className="progress-fill"
          style={{ width: `${percent}%` }}
        />
      </div>
    </div>
  );
}
