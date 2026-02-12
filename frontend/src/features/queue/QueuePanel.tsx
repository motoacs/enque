import type { ChangeEvent } from 'react';
import { useEditStore } from '../../stores/editStore';
import { useEncodeStore } from '../../stores/encodeStore';
import { t } from '../../lib/i18n';

function newJobId(): string {
  return Math.random().toString(36).slice(2, 10);
}

export function QueuePanel() {
  const jobs = useEditStore((s) => s.jobs);
  const queueLocked = useEditStore((s) => s.queueLocked);
  const addJobs = useEditStore((s) => s.addJobs);
  const removeJob = useEditStore((s) => s.removeJob);
  const runtimeJobs = useEncodeStore((s) => s.jobs);

  const onFileSelect = async (e: ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files || files.length === 0) {
      return;
    }
    const normalized = Array.from(files).map((f) => {
      const file = f as File & { path?: string };
      return {
        job_id: newJobId(),
        input_path: file.path ?? f.name,
        input_size_bytes: f.size,
        status: 'pending' as const,
        progress: { percent: 0, fps: null, bitrate_kbps: null, eta_sec: null },
        error_message: ''
      };
    });
    addJobs(normalized);
    e.currentTarget.value = '';
  };

  return (
    <section className="panel">
      <div className="panel-header">
        <h2>{t('queue.title')}</h2>
        <label className="button disabled-when-locked">
          + Files
          <input type="file" multiple disabled={queueLocked} onChange={onFileSelect} />
        </label>
      </div>
      <div className="table">
        {jobs.map((job) => {
          const runtime = runtimeJobs[job.job_id];
          return (
            <div className="row" key={job.job_id}>
              <div className="cell path">{job.input_path}</div>
              <div className="cell status">{runtime?.status ?? job.status}</div>
              <div className="cell progress">{runtime?.percent != null ? `${runtime.percent.toFixed(1)}%` : '-'}</div>
              <button disabled={queueLocked} onClick={() => removeJob(job.job_id)}>
                remove
              </button>
            </div>
          );
        })}
      </div>
    </section>
  );
}
