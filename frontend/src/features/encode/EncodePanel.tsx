import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useEncodeStore } from "@/stores/encodeStore";
import { OverallProgress } from "./OverallProgress";
import { JobProgressList } from "./JobProgressList";
import { LogViewer } from "./LogViewer";
import { EncodeControls } from "./EncodeControls";
import { OverwriteDialog } from "./OverwriteDialog";
import { SessionSummary } from "./SessionSummary";
import * as api from "@/lib/api";

export function EncodePanel() {
  const { t } = useTranslation();
  const [selectedJobId, setSelectedJobId] = useState("");
  const sessionId = useEncodeStore((s) => s.sessionId);
  const sessionState = useEncodeStore((s) => s.sessionState);
  const sessionSummary = useEncodeStore((s) => s.sessionSummary);
  const overwriteRequest = useEncodeStore((s) => s.overwriteRequest);
  const resetSession = useEncodeStore((s) => s.resetSession);
  const clearOverwriteRequest = useEncodeStore((s) => s.clearOverwriteRequest);

  const handleStop = async () => {
    if (sessionId) {
      await api.requestGracefulStop(sessionId);
    }
  };

  const handleAbort = async () => {
    if (sessionId) {
      await api.requestAbort(sessionId);
    }
  };

  const handleOverwrite = async () => {
    if (overwriteRequest) {
      await api.resolveOverwrite(overwriteRequest.sessionId, overwriteRequest.jobId, "overwrite");
      clearOverwriteRequest();
    }
  };

  const handleSkip = async () => {
    if (overwriteRequest) {
      await api.resolveOverwrite(overwriteRequest.sessionId, overwriteRequest.jobId, "skip");
      clearOverwriteRequest();
    }
  };

  const handleOverwriteAbort = async () => {
    if (overwriteRequest) {
      await api.resolveOverwrite(overwriteRequest.sessionId, overwriteRequest.jobId, "abort");
      clearOverwriteRequest();
    }
  };

  const handleDismissSummary = () => {
    resetSession();
  };

  return (
    <div className="flex flex-col h-full">
      <OverallProgress />

      <div className="flex-1 flex overflow-hidden">
        {/* Left: job list */}
        <div className="w-64 flex flex-col overflow-hidden" style={{ borderRight: '1px solid rgba(255,255,255,0.06)', background: 'rgba(18, 18, 26, 0.3)' }}>
          <div className="px-4 py-2.5" style={{ borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
            <h3 className="font-display text-xs font-semibold uppercase tracking-wider" style={{ color: '#9d9da7' }}>
              {t("encode.jobs")}
            </h3>
          </div>
          <JobProgressList
            selectedJobId={selectedJobId}
            onSelectJob={setSelectedJobId}
          />
        </div>

        {/* Right: log viewer */}
        <div className="flex-1 flex flex-col overflow-hidden">
          <div className="px-4 py-2.5" style={{ borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
            <h3 className="font-display text-xs font-semibold uppercase tracking-wider" style={{ color: '#9d9da7' }}>
              {t("encode.log")}
            </h3>
          </div>
          {selectedJobId ? (
            <LogViewer jobId={selectedJobId} />
          ) : (
            <div className="flex-1 flex items-center justify-center text-xs" style={{ color: '#3c3c48' }}>
              {t("encode.selectJobForLog")}
            </div>
          )}
        </div>
      </div>

      {(sessionState === "running" || sessionState === "stopping" || sessionState === "aborting") && (
        <EncodeControls onStop={handleStop} onAbort={handleAbort} />
      )}

      {overwriteRequest && (
        <OverwriteDialog
          outputPath={overwriteRequest.outputPath}
          onOverwrite={handleOverwrite}
          onSkip={handleSkip}
          onAbort={handleOverwriteAbort}
        />
      )}

      {sessionSummary && (sessionState === "completed" || sessionState === "aborted") && (
        <SessionSummary summary={sessionSummary} onDismiss={handleDismissSummary} />
      )}
    </div>
  );
}
