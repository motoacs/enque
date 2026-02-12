import { useTranslation } from "react-i18next";
import { X } from "lucide-react";

interface GPUInfoDialogProps {
  open: boolean;
  onClose: () => void;
  gpuInfo: string;
}

export function GPUInfoDialog({ open, onClose, gpuInfo }: GPUInfoDialogProps) {
  const { t } = useTranslation();

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-zinc-800 rounded-lg shadow-xl w-[600px] max-h-[80vh] flex flex-col border border-zinc-700">
        <div className="flex items-center justify-between px-4 py-3 border-b border-zinc-700">
          <h2 className="text-sm font-semibold text-zinc-200">{t("settings.gpuInfo")}</h2>
          <button onClick={onClose} className="text-zinc-400 hover:text-zinc-200">
            <X size={16} />
          </button>
        </div>
        <div className="flex-1 overflow-y-auto p-4">
          <pre className="text-xs text-zinc-300 font-mono whitespace-pre-wrap">
            {gpuInfo || t("settings.noGpuInfo")}
          </pre>
        </div>
        <div className="px-4 py-3 border-t border-zinc-700 flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-1.5 text-xs bg-zinc-600 hover:bg-zinc-500 text-zinc-200 rounded"
          >
            {t("common.close")}
          </button>
        </div>
      </div>
    </div>
  );
}
