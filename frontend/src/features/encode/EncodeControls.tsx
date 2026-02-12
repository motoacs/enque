import { useTranslation } from "react-i18next";
import { useEncodeStore, type SessionState } from "@/stores/encodeStore";
import { Square, StopCircle, Play } from "lucide-react";

interface EncodeControlsProps {
  onStop: () => void;
  onAbort: () => void;
}

export function EncodeControls({ onStop, onAbort }: EncodeControlsProps) {
  const { t } = useTranslation();
  const sessionState = useEncodeStore((s) => s.sessionState);

  return (
    <div className="flex items-center gap-2 px-4 py-2 border-t border-zinc-700 bg-zinc-800">
      {sessionState === "running" && (
        <>
          <button
            onClick={onStop}
            className="flex items-center gap-1.5 px-3 py-1.5 text-xs bg-yellow-600 hover:bg-yellow-500 text-white rounded transition-colors"
          >
            <StopCircle size={14} />
            {t("encode.stop")}
          </button>
          <button
            onClick={onAbort}
            className="flex items-center gap-1.5 px-3 py-1.5 text-xs bg-red-600 hover:bg-red-500 text-white rounded transition-colors"
          >
            <Square size={14} />
            {t("encode.abort")}
          </button>
        </>
      )}
      {sessionState === "stopping" && (
        <span className="text-xs text-yellow-400">
          {t("encode.stoppingMsg")}
        </span>
      )}
      {sessionState === "aborting" && (
        <span className="text-xs text-red-400">
          {t("encode.abortingMsg")}
        </span>
      )}
    </div>
  );
}
