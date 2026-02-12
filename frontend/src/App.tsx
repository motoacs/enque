import { useEffect, useMemo, useRef, useState } from 'react';
import { QueuePanel } from './features/queue/QueuePanel';
import { ProfilePanel } from './features/profile/ProfilePanel';
import { SettingsPanel } from './features/settings/SettingsPanel';
import { EncodeMonitor } from './features/encode/EncodeMonitor';
import { api } from './lib/api';
import { onEvent } from './lib/events';
import { t } from './lib/i18n';
import { useAppStore } from './stores/appStore';
import { useEditStore } from './stores/editStore';
import { useEncodeStore } from './stores/encodeStore';
import { useProfileStore } from './stores/profileStore';

export default function App() {
  const [preview, setPreview] = useState('');
  const overwriteDecisionForAll = useRef<'overwrite' | 'skip' | 'abort' | null>(null);
  const jobs = useEditStore((s) => s.jobs);
  const setQueueLocked = useEditStore((s) => s.setQueueLocked);
  const config = useAppStore((s) => s.config);
  const setConfig = useAppStore((s) => s.setConfig);
  const setTools = useAppStore((s) => s.setTools);
  const setGPUInfo = useAppStore((s) => s.setGPUInfo);
  const setWarnings = useAppStore((s) => s.setWarnings);
  const profiles = useProfileStore((s) => s.profiles);
  const selectedProfileId = useProfileStore((s) => s.selectedProfileId);
  const selectedProfile = useMemo(() => profiles.find((p) => p.id === selectedProfileId), [profiles, selectedProfileId]);
  const setProfiles = useProfileStore((s) => s.setProfiles);

  const encodeSessionId = useEncodeStore((s) => s.sessionId);
  const encodeState = useEncodeStore((s) => s.state);
  const setSession = useEncodeStore((s) => s.setSession);
  const setEncodeState = useEncodeStore((s) => s.setState);
  const setCompletedJobs = useEncodeStore((s) => s.setCompletedJobs);
  const upsertJob = useEncodeStore((s) => s.upsertJob);
  const appendJobLog = useEncodeStore((s) => s.appendJobLog);

  useEffect(() => {
    api
      .bootstrap()
      .then((res) => {
        setConfig(res.config);
        setProfiles(res.profiles);
        setTools(res.tools);
        setGPUInfo(res.gpu_info);
        setWarnings(res.warnings);
      })
      .catch(() => {
        // Dev mode without Wails runtime.
      });
  }, [setConfig, setGPUInfo, setProfiles, setTools, setWarnings]);

  useEffect(() => {
    const unsubscribers = [
      onEvent<{ session_id: string; total_jobs: number }>('enque:session_started', (p) => {
        setSession(p.session_id, p.total_jobs);
        setQueueLocked(true);
      }),
      onEvent<{ session_id: string; completed_jobs: number; stop_requested?: boolean; abort_requested?: boolean }>('enque:session_state', (p) => {
        setCompletedJobs(p.completed_jobs);
        if (p.abort_requested) setEncodeState('aborting');
        else if (p.stop_requested) setEncodeState('stopping');
      }),
      onEvent<{ session_id: string }>('enque:session_finished', () => {
        setEncodeState('completed');
        setQueueLocked(false);
        overwriteDecisionForAll.current = null;
      }),
      onEvent<{ session_id: string; job_id: string; final_output_path: string }>('enque:job_needs_overwrite', async (p) => {
        let decision = overwriteDecisionForAll.current;
        if (!decision) {
          const overwrite = window.confirm(`Output already exists:\\n${p.final_output_path}\\n\\nOverwrite this file?\\nCancel = skip`);
          decision = overwrite ? 'overwrite' : 'skip';
          const applyAll = window.confirm('Apply this decision to the rest of this session?');
          if (applyAll) {
            overwriteDecisionForAll.current = decision;
          }
        }
        await api.resolveOverwrite(p.session_id, p.job_id, decision);
      }),
      onEvent<{ job_id: string; status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled' | 'timeout' | 'skipped'; exit_code?: number; error_message?: string }>(
        'enque:job_finished',
        (p) => {
          upsertJob(p.job_id, { status: p.status, exit_code: p.exit_code, error_message: p.error_message });
        }
      ),
      onEvent<{ job_id: string; percent: number | null; fps: number | null; bitrate_kbps: number | null; eta_sec: number | null }>(
        'enque:job_progress',
        (p) => {
          upsertJob(p.job_id, {
            percent: p.percent,
            fps: p.fps,
            bitrate_kbps: p.bitrate_kbps,
            eta_sec: p.eta_sec
          });
        }
      ),
      onEvent<{ job_id: string; line: string }>('enque:job_log', (p) => {
        appendJobLog(p.job_id, p.line);
      })
    ];
    return () => {
      unsubscribers.forEach((unsub) => unsub());
    };
  }, [appendJobLog, setCompletedJobs, setEncodeState, setQueueLocked, setSession, upsertJob]);

  useEffect(() => {
    if (!selectedProfile || jobs.length === 0) {
      setPreview('');
      return;
    }
    const input = jobs[0]?.input_path ?? 'input.mp4';
    const output = `${input}.preview.${config.output_container}`;
    api
      .previewCommand({
        profile: selectedProfile,
        app_config_snapshot: config,
        input_path: input,
        output_path: output
      })
      .then((res) => setPreview(res.display_command))
      .catch(() => {
        setPreview('');
      });
  }, [config, jobs, selectedProfile]);

  const handleStart = async () => {
    if (!selectedProfile || jobs.length === 0) return;
    await api.startEncode({
      jobs: jobs.map((j) => ({ job_id: j.job_id, input_path: j.input_path })),
      profile: selectedProfile,
      app_config_snapshot: config,
      command_preview: preview
    });
  };

  return (
    <main className="app-root">
      <header className="hero">
        <h1>{t('app.title')}</h1>
        <div className="actions">
          <button onClick={handleStart}>{t('button.start')}</button>
          <button disabled={!encodeSessionId} onClick={() => api.requestGracefulStop(encodeSessionId)}>
            {t('button.stop')}
          </button>
          <button disabled={!encodeSessionId} onClick={() => api.requestAbort(encodeSessionId)}>
            {t('button.abort')}
          </button>
          <span className="state-tag">{encodeState}</span>
        </div>
      </header>

      <section className="preview-panel">
        <h2>{t('preview.title')}</h2>
        <textarea readOnly rows={3} value={preview} />
      </section>

      <div className="layout-grid">
        <QueuePanel />
        <ProfilePanel />
        <SettingsPanel />
        <EncodeMonitor />
      </div>
    </main>
  );
}
