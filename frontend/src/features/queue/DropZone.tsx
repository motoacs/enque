import { useCallback, useEffect, useRef, useState, type ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { Film } from "lucide-react";
import { OnFileDrop, OnFileDropOff } from "../../../wailsjs/runtime/runtime";

interface DropZoneProps {
  onFilesDropped: (paths: string[]) => void;
  children: ReactNode;
}

export function DropZone({ onFilesDropped, children }: DropZoneProps) {
  const { t } = useTranslation();
  const [isDragging, setIsDragging] = useState(false);
  const dragCounter = useRef(0);
  const callbackRef = useRef(onFilesDropped);
  callbackRef.current = onFilesDropped;

  // Register Wails native file drop handler for actual file paths
  useEffect(() => {
    OnFileDrop((_x: number, _y: number, paths: string[]) => {
      dragCounter.current = 0;
      setIsDragging(false);
      if (paths && paths.length > 0) {
        callbackRef.current(paths);
      }
    }, true);

    return () => {
      OnFileDropOff();
    };
  }, []);

  // Visual feedback via HTML5 drag events with counter to prevent flicker
  const handleDragEnter = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    dragCounter.current++;
    if (dragCounter.current === 1) {
      setIsDragging(true);
    }
  }, []);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    dragCounter.current--;
    if (dragCounter.current <= 0) {
      dragCounter.current = 0;
      setIsDragging(false);
    }
  }, []);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    dragCounter.current = 0;
    setIsDragging(false);
    // Actual file handling is done by Wails OnFileDrop callback
  }, []);

  return (
    <div
      onDragEnter={handleDragEnter}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      className="relative h-full"
      style={{ "--wails-drop-target": "drop" } as React.CSSProperties}
    >
      {isDragging && (
        <div
          className="absolute inset-1 z-10 flex flex-col items-center justify-center gap-3 rounded-xl animate-fade-in pointer-events-none"
          style={{
            background: 'rgba(232, 168, 73, 0.06)',
            border: '2px dashed rgba(232, 168, 73, 0.35)',
          }}
        >
          <Film size={28} style={{ color: '#e8a849', opacity: 0.7 }} />
          <p className="font-display text-sm font-semibold" style={{ color: '#e8a849' }}>
            {t("queue.dropHere")}
          </p>
        </div>
      )}
      {children}
    </div>
  );
}
