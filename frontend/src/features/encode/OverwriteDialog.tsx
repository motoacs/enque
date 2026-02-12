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
    <div className="dialog-overlay">
      <div className="dialog-panel w-[480px]">
        <div className="dialog-header">
          <div className="flex items-center gap-2.5">
            <AlertTriangle size={16} style={{ color: '#fbbf24' }} />
            <h2 className="text-sm font-display font-semibold" style={{ color: '#e8e6e3' }}>
              {t("encode.overwriteTitle")}
            </h2>
          </div>
        </div>
        <div className="p-5">
          <p className="text-xs mb-3" style={{ color: '#9d9da7' }}>{t("encode.overwriteMsg")}</p>
          <p
            className="text-xs font-mono rounded-lg p-3 break-all"
            style={{
              background: 'rgba(10, 10, 15, 0.8)',
              color: '#9d9da7',
              border: '1px solid rgba(255,255,255,0.04)',
            }}
          >
            {outputPath}
          </p>
        </div>
        <div className="dialog-footer">
          <button onClick={onAbort} className="btn-danger">
            {t("encode.abort")}
          </button>
          <button onClick={onSkip} className="btn-secondary">
            {t("encode.skipFile")}
          </button>
          <button onClick={onOverwrite} className="btn-warning">
            {t("encode.overwrite")}
          </button>
        </div>
      </div>
    </div>
  );
}
