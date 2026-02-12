import { useTranslation } from "react-i18next";
import { AlertTriangle } from "lucide-react";

interface OverwriteDialogProps {
  outputPath: string;
  onOverwrite: () => void;
  onSkip: () => void;
  onAbort: () => void;
}

export function OverwriteDialog({ outputPath, onOverwrite, onSkip, onAbort }: OverwriteDialogProps) {
  const { t } = useTranslation();

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-zinc-800 rounded-lg shadow-xl w-[450px] border border-zinc-700">
        <div className="flex items-center gap-2 px-4 py-3 border-b border-zinc-700">
          <AlertTriangle size={16} className="text-yellow-400" />
          <h2 className="text-sm font-semibold text-zinc-200">{t("encode.overwriteTitle")}</h2>
        </div>
        <div className="p-4">
          <p className="text-xs text-zinc-300 mb-2">{t("encode.overwriteMsg")}</p>
          <p className="text-xs text-zinc-400 font-mono bg-zinc-900 p-2 rounded break-all">
            {outputPath}
          </p>
        </div>
        <div className="px-4 py-3 border-t border-zinc-700 flex justify-end gap-2">
          <button
            onClick={onAbort}
            className="px-3 py-1.5 text-xs bg-red-600 hover:bg-red-500 text-white rounded transition-colors"
          >
            {t("encode.abort")}
          </button>
          <button
            onClick={onSkip}
            className="px-3 py-1.5 text-xs bg-zinc-600 hover:bg-zinc-500 text-zinc-200 rounded transition-colors"
          >
            {t("encode.skipFile")}
          </button>
          <button
            onClick={onOverwrite}
            className="px-3 py-1.5 text-xs bg-yellow-600 hover:bg-yellow-500 text-white rounded transition-colors"
          >
            {t("encode.overwrite")}
          </button>
        </div>
      </div>
    </div>
  );
}
