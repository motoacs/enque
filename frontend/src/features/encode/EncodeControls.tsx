import { useTranslation } from "react-i18next";
import { useEncodeStore, type SessionState } from "@/stores/encodeStore";
import { Square, StopCircle } from "lucide-react";

interface EncodeControlsProps {
  onStop: () => void;
  onAbort: () => void;
}

export function EncodeControls({ onStop, onAbort }: EncodeControlsProps) {
  const { t } = useTranslation();
  const sessionState = useEncodeStore((s) => s.sessionState);

  return (
    <div className="flex items-center gap-3 px-5 py-3" style={{ borderTop: '1px solid rgba(255,255,255,0.06)', background: 'rgba(18, 18, 26, 0.5)' }}>
      {sessionState === "running" && (
        <>
          <button onClick={onStop} className="btn-warning">
            <StopCircle size={13} />
            {t("encode.stop")}
          </button>
          <button onClick={onAbort} className="btn-danger">
            <Square size={13} />
            {t("encode.abort")}
          </button>
        </>
      )}
      {sessionState === "stopping" && (
        <span className="text-xs font-display font-medium" style={{ color: '#fbbf24' }}>
          {t("encode.stoppingMsg")}
        </span>
      )}
      {sessionState === "aborting" && (
        <span className="text-xs font-display font-medium" style={{ color: '#f87171' }}>
          {t("encode.abortingMsg")}
        </span>
      )}
    </div>
  );
}
