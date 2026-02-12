import { useRef, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useEncodeStore } from "@/stores/encodeStore";

interface LogViewerProps {
  jobId: string;
}

export function LogViewer({ jobId }: LogViewerProps) {
  const { t } = useTranslation();
  const logs = useEncodeStore((s) => s.jobLogs[jobId] || []);
  const containerRef = useRef<HTMLDivElement>(null);
  const autoScrollRef = useRef(true);

  useEffect(() => {
    if (autoScrollRef.current && containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
    }
  }, [logs.length]);

  const handleScroll = () => {
    if (!containerRef.current) return;
    const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
    autoScrollRef.current = scrollHeight - scrollTop - clientHeight < 40;
  };

  return (
    <div
      ref={containerRef}
      onScroll={handleScroll}
      className="flex-1 overflow-y-auto p-3 font-mono text-[11px] leading-relaxed"
      style={{
        background: 'rgba(8, 8, 12, 0.6)',
        color: '#5c5c68',
      }}
    >
      {logs.length === 0 && (
        <div className="text-center py-8" style={{ color: '#3c3c48' }}>{t("encode.selectJobForLog")}</div>
      )}
      {logs.map((line, i) => (
        <div key={i} className="whitespace-pre-wrap break-all hover:bg-white/[0.02] px-1 rounded">
          {line}
        </div>
      ))}
    </div>
  );
}
