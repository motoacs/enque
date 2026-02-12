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
    <div className="dialog-overlay">
      <div className="dialog-panel w-[520px] max-h-[60vh]">
        <div className="dialog-header">
          <div className="flex items-center gap-2.5">
            <Trash2 size={15} style={{ color: '#fbbf24' }} />
            <h2 className="text-sm font-display font-semibold" style={{ color: '#e8e6e3' }}>
              {t("settings.tempCleanup")}
            </h2>
          </div>
          <button onClick={onDismiss} className="icon-btn">
            <X size={15} />
          </button>
        </div>
        <div className="flex-1 overflow-y-auto p-5">
          <p className="text-xs mb-3" style={{ color: '#9d9da7' }}>{t("settings.tempCleanupMsg")}</p>
          <div className="space-y-1.5">
            {tempFiles.map((path, i) => (
              <div
                key={i}
                className="text-xs font-mono rounded-md px-3 py-1.5 truncate"
                style={{
                  background: 'rgba(10, 10, 15, 0.8)',
                  color: '#9d9da7',
                  border: '1px solid rgba(255,255,255,0.04)',
                }}
                title={path}
              >
                {path}
              </div>
            ))}
          </div>
        </div>
        <div className="dialog-footer">
          <button onClick={onDismiss} className="btn-secondary">
            {t("settings.keepFiles")}
          </button>
          <button onClick={onCleanup} className="btn-danger">
            {t("settings.deleteAll")}
          </button>
        </div>
      </div>
    </div>
  );
}
