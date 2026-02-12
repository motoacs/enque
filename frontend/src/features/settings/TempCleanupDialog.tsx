import { useTranslation } from "react-i18next";
import { Trash2, X } from "lucide-react";

interface TempCleanupDialogProps {
  open: boolean;
  tempFiles: string[];
  onCleanup: () => void;
  onDismiss: () => void;
}

export function TempCleanupDialog({ open, tempFiles, onCleanup, onDismiss }: TempCleanupDialogProps) {
  const { t } = useTranslation();

  if (!open || tempFiles.length === 0) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-zinc-800 rounded-lg shadow-xl w-[500px] max-h-[60vh] flex flex-col border border-zinc-700">
        <div className="flex items-center justify-between px-4 py-3 border-b border-zinc-700">
          <div className="flex items-center gap-2">
            <Trash2 size={16} className="text-yellow-400" />
            <h2 className="text-sm font-semibold text-zinc-200">{t("settings.tempCleanup")}</h2>
          </div>
          <button onClick={onDismiss} className="text-zinc-400 hover:text-zinc-200">
            <X size={16} />
          </button>
        </div>
        <div className="flex-1 overflow-y-auto p-4">
          <p className="text-xs text-zinc-400 mb-3">{t("settings.tempCleanupMsg")}</p>
          <div className="space-y-1">
            {tempFiles.map((path, i) => (
              <div key={i} className="text-xs text-zinc-300 font-mono bg-zinc-900 px-2 py-1 rounded truncate" title={path}>
                {path}
              </div>
            ))}
          </div>
        </div>
        <div className="px-4 py-3 border-t border-zinc-700 flex justify-end gap-2">
          <button
            onClick={onDismiss}
            className="px-4 py-1.5 text-xs bg-zinc-600 hover:bg-zinc-500 text-zinc-200 rounded transition-colors"
          >
            {t("settings.keepFiles")}
          </button>
          <button
            onClick={onCleanup}
            className="px-4 py-1.5 text-xs bg-red-600 hover:bg-red-500 text-white rounded transition-colors"
          >
            {t("settings.deleteAll")}
          </button>
        </div>
      </div>
    </div>
  );
}
