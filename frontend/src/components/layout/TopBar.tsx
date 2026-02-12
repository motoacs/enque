import { useTranslation } from "react-i18next";
import { Settings } from "lucide-react";

interface TopBarProps {
  onSettingsClick: () => void;
}

export function TopBar({ onSettingsClick }: TopBarProps) {
  const { t } = useTranslation();

  return (
    <header className="flex items-center justify-between px-5 py-2.5 shrink-0 glass-panel wails-drag" style={{ borderTop: 'none', borderLeft: 'none', borderRight: 'none' }}>
      <div className="flex items-center gap-3">
        {/* Amber accent mark */}
        <div className="w-2 h-5 rounded-sm" style={{ background: 'linear-gradient(180deg, #e8a849, #d4922a)' }} />
        <h1 className="font-display text-base font-bold tracking-wide" style={{ color: '#e8e6e3' }}>
          {t("app.title")}
        </h1>
      </div>
      <button
        onClick={onSettingsClick}
        className="icon-btn"
      >
        <Settings size={16} />
      </button>
    </header>
  );
}
