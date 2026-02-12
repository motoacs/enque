import { useTranslation } from "react-i18next";
import { AlertTriangle } from "lucide-react";

interface DeleteConfirmDialogProps {
  open: boolean;
  profileName: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export function DeleteConfirmDialog({ open, profileName, onConfirm, onCancel }: DeleteConfirmDialogProps) {
  const { t } = useTranslation();

  if (!open) return null;

  return (
    <div className="dialog-overlay">
      <div className="dialog-panel w-[400px]">
        <div className="dialog-header">
          <div className="flex items-center gap-2.5">
            <AlertTriangle size={16} style={{ color: '#fbbf24' }} />
            <h2 className="text-sm font-display font-semibold" style={{ color: '#e8e6e3' }}>
              {t("profile.delete")}
            </h2>
          </div>
        </div>
        <div className="p-5">
          <p className="text-xs" style={{ color: '#9d9da7' }}>
            {t("profile.deleteConfirm", { name: profileName })}
          </p>
        </div>
        <div className="dialog-footer">
          <button onClick={onCancel} className="btn-secondary">
            {t("common.cancel")}
          </button>
          <button onClick={onConfirm} className="btn-danger">
            {t("profile.delete")}
          </button>
        </div>
      </div>
    </div>
  );
}
