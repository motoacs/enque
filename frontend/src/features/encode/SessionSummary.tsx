import { useTranslation } from "react-i18next";
import type { SessionSummary as SessionSummaryType } from "@/stores/encodeStore";
import { CheckCircle, XCircle, MinusCircle, AlertTriangle, Clock } from "lucide-react";

interface SessionSummaryProps {
  summary: SessionSummaryType;
  onDismiss: () => void;
}

export function SessionSummary({ summary, onDismiss }: SessionSummaryProps) {
  const { t } = useTranslation();

  const isSuccess = summary.failedJobs === 0 && summary.timeoutJobs === 0;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-zinc-800 rounded-lg shadow-xl w-[400px] border border-zinc-700">
        <div className="flex items-center gap-2 px-4 py-3 border-b border-zinc-700">
          {isSuccess ? (
            <CheckCircle size={16} className="text-green-400" />
          ) : (
            <AlertTriangle size={16} className="text-yellow-400" />
          )}
          <h2 className="text-sm font-semibold text-zinc-200">
            {t("encode.sessionComplete")}
          </h2>
        </div>

        <div className="p-4 space-y-2">
          <div className="grid grid-cols-2 gap-2 text-xs">
            <div className="flex items-center gap-1.5 text-zinc-400">
              <span>{t("encode.totalJobs")}:</span>
            </div>
            <div className="text-zinc-200">{summary.totalJobs}</div>

            <div className="flex items-center gap-1.5 text-green-400">
              <CheckCircle size={12} />
              <span>{t("encode.completed")}:</span>
            </div>
            <div className="text-zinc-200">{summary.completedJobs}</div>

            {summary.failedJobs > 0 && (
              <>
                <div className="flex items-center gap-1.5 text-red-400">
                  <XCircle size={12} />
                  <span>{t("encode.failed")}:</span>
                </div>
                <div className="text-zinc-200">{summary.failedJobs}</div>
              </>
            )}

            {summary.cancelledJobs > 0 && (
              <>
                <div className="flex items-center gap-1.5 text-yellow-400">
                  <MinusCircle size={12} />
                  <span>{t("encode.cancelled")}:</span>
                </div>
                <div className="text-zinc-200">{summary.cancelledJobs}</div>
              </>
            )}

            {summary.timeoutJobs > 0 && (
              <>
                <div className="flex items-center gap-1.5 text-orange-400">
                  <Clock size={12} />
                  <span>{t("encode.timedOut")}:</span>
                </div>
                <div className="text-zinc-200">{summary.timeoutJobs}</div>
              </>
            )}

            {summary.skippedJobs > 0 && (
              <>
                <div className="flex items-center gap-1.5 text-zinc-500">
                  <MinusCircle size={12} />
                  <span>{t("encode.skipped")}:</span>
                </div>
                <div className="text-zinc-200">{summary.skippedJobs}</div>
              </>
            )}
          </div>
        </div>

        <div className="px-4 py-3 border-t border-zinc-700 flex justify-end">
          <button
            onClick={onDismiss}
            className="px-4 py-1.5 text-xs bg-blue-600 hover:bg-blue-500 text-white rounded transition-colors"
          >
            {t("common.close")}
          </button>
        </div>
      </div>
    </div>
  );
}
