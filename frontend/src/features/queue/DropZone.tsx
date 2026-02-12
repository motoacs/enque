import { useCallback, useState, type DragEvent, type ReactNode } from "react";
import { useTranslation } from "react-i18next";

interface DropZoneProps {
  onFilesDropped: (paths: string[]) => void;
  children: ReactNode;
}

export function DropZone({ onFilesDropped, children }: DropZoneProps) {
  const { t } = useTranslation();
  const [isDragging, setIsDragging] = useState(false);

  const handleDragOver = useCallback((e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback(
    (e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setIsDragging(false);

      const files = e.dataTransfer?.files;
      if (files && files.length > 0) {
        const paths: string[] = [];
        for (let i = 0; i < files.length; i++) {
          // Wails provides the full path via webkitRelativePath or name
          const f = files[i] as File & { path?: string };
          paths.push(f.path || f.name);
        }
        onFilesDropped(paths);
      }
    },
    [onFilesDropped]
  );

  return (
    <div
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      className="relative h-full"
    >
      {isDragging && (
        <div className="absolute inset-0 z-10 flex items-center justify-center bg-blue-600/20 border-2 border-dashed border-blue-500 rounded-lg">
          <p className="text-blue-400 font-medium text-lg">
            {t("queue.dropHere")}
          </p>
        </div>
      )}
      {children}
    </div>
  );
}
