import { useTranslation } from "react-i18next";
import { useState } from "react";
import { Copy, Check } from "lucide-react";

interface CommandPreviewProps {
  command: string;
}

export function CommandPreview({ command }: CommandPreviewProps) {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(command);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {}
  };

  return (
    <section className="px-3 py-2 border-t border-zinc-700">
      <div className="flex items-center justify-between mb-1">
        <h3 className="text-xs font-semibold text-zinc-400 uppercase">
          {t("preview.title")}
        </h3>
        <button
          onClick={handleCopy}
          className="flex items-center gap-1 text-xs text-zinc-400 hover:text-zinc-200 transition-colors"
          disabled={!command}
        >
          {copied ? <Check size={12} /> : <Copy size={12} />}
          {copied ? t("preview.copied") : t("preview.copy")}
        </button>
      </div>
      <pre className="text-xs text-zinc-300 bg-zinc-800 rounded p-2 overflow-x-auto font-mono whitespace-pre-wrap break-all max-h-24">
        {command || "..."}
      </pre>
    </section>
  );
}
