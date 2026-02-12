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
    <div className="dialog-overlay">
      <div className="dialog-panel w-[420px]">
        <div className="dialog-header">
          <div className="flex items-center gap-2.5">
            {isSuccess ? (
              <CheckCircle size={16} style={{ color: '#34d399' }} />
            ) : (
              <AlertTriangle size={16} style={{ color: '#fbbf24' }} />
            )}
            <h2 className="text-sm font-display font-semibold" style={{ color: '#e8e6e3' }}>
              {t("encode.sessionComplete")}
            </h2>
          </div>
        </div>

        <div className="p-5 space-y-2.5">
          <div className="grid grid-cols-2 gap-2 text-xs">
            <div style={{ color: '#9d9da7' }}>
              {t("encode.totalJobs")}:
            </div>
            <div className="font-mono" style={{ color: '#e8e6e3' }}>{summary.totalJobs}</div>

            <div className="flex items-center gap-1.5" style={{ color: '#34d399' }}>
              <CheckCircle size={11} />
              {t("encode.completed")}:
            </div>
            <div className="font-mono" style={{ color: '#e8e6e3' }}>{summary.completedJobs}</div>

            {summary.failedJobs > 0 && (
              <>
                <div className="flex items-center gap-1.5" style={{ color: '#f87171' }}>
                  <XCircle size={11} />
                  {t("encode.failed")}:
                </div>
                <div className="font-mono" style={{ color: '#e8e6e3' }}>{summary.failedJobs}</div>
              </>
            )}

            {summary.cancelledJobs > 0 && (
              <>
                <div className="flex items-center gap-1.5" style={{ color: '#fbbf24' }}>
                  <MinusCircle size={11} />
                  {t("encode.cancelled")}:
                </div>
                <div className="font-mono" style={{ color: '#e8e6e3' }}>{summary.cancelledJobs}</div>
              </>
            )}

            {summary.timeoutJobs > 0 && (
              <>
                <div className="flex items-center gap-1.5" style={{ color: '#fb923c' }}>
                  <Clock size={11} />
                  {t("encode.timedOut")}:
                </div>
                <div className="font-mono" style={{ color: '#e8e6e3' }}>{summary.timeoutJobs}</div>
              </>
            )}

            {summary.skippedJobs > 0 && (
              <>
                <div className="flex items-center gap-1.5" style={{ color: '#5c5c68' }}>
                  <MinusCircle size={11} />
                  {t("encode.skipped")}:
                </div>
                <div className="font-mono" style={{ color: '#e8e6e3' }}>{summary.skippedJobs}</div>
              </>
            )}
          </div>
        </div>

        <div className="dialog-footer">
          <button onClick={onDismiss} className="btn-primary">
            {t("common.close")}
          </button>
        </div>
      </div>
    </div>
  );
}
