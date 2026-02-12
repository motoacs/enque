import { useMemo, useState } from 'react';
import { useEncodeStore } from '../../stores/encodeStore';
import { t } from '../../lib/i18n';

export function EncodeMonitor() {
  const jobs = useEncodeStore((s) => s.jobs);
  const state = useEncodeStore((s) => s.state);
  const completed = useEncodeStore((s) => s.completedJobs);
  const total = useEncodeStore((s) => s.totalJobs);
  const [selectedJobID, setSelectedJobID] = useState<string>('');

  const jobIDs = Object.keys(jobs);
  const selected = useMemo(() => jobs[selectedJobID] ?? jobs[jobIDs[0]], [jobs, jobIDs, selectedJobID]);

  return (
    <section className="panel">
      <div className="panel-header">
        <h2>{t('encode.title')}</h2>
        <div>
          {state} ({completed}/{total})
        </div>
      </div>
      <div className="monitor-grid">
        <div>
          {jobIDs.map((id) => (
            <button key={id} className="job-pill" onClick={() => setSelectedJobID(id)}>
              {id}: {jobs[id].status}
            </button>
          ))}
        </div>
        <pre className="log-box">{selected ? selected.logs.join('\n') : 'No logs yet'}</pre>
      </div>
    </section>
  );
}
