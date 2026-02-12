import { useState, useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";

interface ProfileNameDialogProps {
  open: boolean;
  title: string;
  initialName: string;
  onConfirm: (name: string) => void;
  onCancel: () => void;
}

export function ProfileNameDialog({ open, title, initialName, onConfirm, onCancel }: ProfileNameDialogProps) {
  const { t } = useTranslation();
  const [name, setName] = useState(initialName);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (open) {
      setName(initialName);
      setTimeout(() => inputRef.current?.select(), 0);
    }
  }, [open, initialName]);

  if (!open) return null;

  const isValid = name.trim().length >= 1 && name.trim().length <= 80;

  const handleConfirm = () => {
    if (isValid) onConfirm(name.trim());
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && isValid) {
      handleConfirm();
    } else if (e.key === "Escape") {
      onCancel();
    }
  };

  return (
    <div className="dialog-overlay">
      <div className="dialog-panel w-[400px]">
        <div className="dialog-header">
          <h2 className="text-sm font-display font-semibold" style={{ color: '#e8e6e3' }}>
            {title}
          </h2>
        </div>
        <div className="p-5">
          <label className="block text-xs mb-2" style={{ color: '#9d9da7' }}>
            {t("profile.nameLabel")}
          </label>
          <input
            ref={inputRef}
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            onKeyDown={handleKeyDown}
            maxLength={80}
            className="w-full form-input"
            autoFocus
          />
        </div>
        <div className="dialog-footer">
          <button onClick={onCancel} className="btn-secondary">
            {t("common.cancel")}
          </button>
          <button onClick={handleConfirm} disabled={!isValid} className="btn-primary">
            {t("common.ok")}
          </button>
        </div>
      </div>
    </div>
  );
}
