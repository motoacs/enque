import { useTranslation } from "react-i18next";
import { Settings } from "lucide-react";

interface TopBarProps {
  onSettingsClick: () => void;
}

export function TopBar({ onSettingsClick }: TopBarProps) {
  const { t } = useTranslation();

  return (
    <header className="flex items-center justify-between px-4 py-2 border-b border-zinc-700 bg-zinc-800 shrink-0">
      <h1 className="text-lg font-semibold text-zinc-100">{t("app.title")}</h1>
      <button
        onClick={onSettingsClick}
        className="p-2 rounded hover:bg-zinc-700 text-zinc-400 hover:text-zinc-200 transition-colors"
      >
        <Settings size={18} />
      </button>
    </header>
  );
}
