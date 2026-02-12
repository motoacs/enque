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
      className="flex-1 overflow-y-auto bg-zinc-900 p-2 font-mono text-[11px] text-zinc-400 leading-tight"
    >
      {logs.length === 0 && (
        <div className="text-zinc-600 text-center py-4">{t("encode.selectJobForLog")}</div>
      )}
      {logs.map((line, i) => (
        <div key={i} className="whitespace-pre-wrap break-all">
          {line}
        </div>
      ))}
    </div>
  );
}
