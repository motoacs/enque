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
    <div className="px-4 py-2 space-y-1">
      <div className="flex justify-between text-xs text-zinc-400">
        <span>{t("encode.progress")}</span>
        <span>
          {completed} / {total} ({percent}%)
        </span>
      </div>
      <div className="h-2 bg-zinc-700 rounded overflow-hidden">
        <div
          className="h-full bg-blue-500 transition-all duration-300"
          style={{ width: `${percent}%` }}
        />
      </div>
    </div>
  );
}
