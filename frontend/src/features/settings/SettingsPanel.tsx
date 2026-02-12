import { useAppStore } from '../../stores/appStore';
import { t } from '../../lib/i18n';

export function SettingsPanel() {
  const config = useAppStore((s) => s.config);
  const patchConfig = useAppStore((s) => s.patchConfig);

  return (
    <section className="panel">
      <div className="panel-header">
        <h2>{t('settings.title')}</h2>
      </div>
      <div className="profile-fields">
        <label>
          Concurrent jobs
          <input
            type="number"
            min={1}
            max={8}
            value={config.max_concurrent_jobs}
            onChange={(e) => patchConfig({ max_concurrent_jobs: Number(e.target.value) })}
          />
        </label>
        <label>
          Overwrite mode
          <select
            value={config.overwrite_mode}
            onChange={(e) => patchConfig({ overwrite_mode: e.target.value as 'ask' | 'auto_rename' })}
          >
            <option value="ask">ask</option>
            <option value="auto_rename">auto_rename</option>
          </select>
        </label>
        <label>
          Output container
          <input value={config.output_container} onChange={(e) => patchConfig({ output_container: e.target.value })} />
        </label>
        <label>
          Output template
          <input value={config.output_name_template} onChange={(e) => patchConfig({ output_name_template: e.target.value })} />
        </label>
      </div>
    </section>
  );
}
