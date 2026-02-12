import { useTranslation } from "react-i18next";
import { useState } from "react";
import { Copy, Check, Terminal } from "lucide-react";

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
    <section className="shrink-0 px-4 py-3" style={{ borderTop: '1px solid rgba(255,255,255,0.08)', background: 'rgba(12, 12, 18, 0.7)' }}>
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <Terminal size={12} style={{ color: '#4ecdc4' }} />
          <h3 className="font-display text-[10px] font-semibold uppercase tracking-wider" style={{ color: '#4ecdc4' }}>
            {t("preview.title")}
          </h3>
        </div>
        <button
          onClick={handleCopy}
          className="flex items-center gap-1.5 text-[10px] transition-colors duration-150"
          style={{ color: copied ? '#34d399' : '#5c5c68' }}
          disabled={!command}
        >
          {copied ? <Check size={10} /> : <Copy size={10} />}
          {copied ? t("preview.copied") : t("preview.copy")}
        </button>
      </div>
      <pre
        className="text-[11px] font-mono rounded-lg p-3 overflow-x-auto whitespace-pre-wrap break-all max-h-24 leading-relaxed"
        style={{
          background: 'rgba(10, 10, 15, 0.8)',
          color: '#9d9da7',
          border: '1px solid rgba(255,255,255,0.04)',
        }}
      >
        {command || "..."}
      </pre>
    </section>
  );
}
