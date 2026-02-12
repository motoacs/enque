import { useMemo } from 'react';
import { useProfileStore } from '../../stores/profileStore';
import { t } from '../../lib/i18n';

export function ProfilePanel() {
  const profiles = useProfileStore((s) => s.profiles);
  const selectedProfileId = useProfileStore((s) => s.selectedProfileId);
  const selectProfile = useProfileStore((s) => s.selectProfile);
  const updateSelectedProfile = useProfileStore((s) => s.updateSelectedProfile);

  const selected = useMemo(() => profiles.find((p) => p.id === selectedProfileId), [profiles, selectedProfileId]);

  return (
    <section className="panel">
      <div className="panel-header">
        <h2>{t('profile.title')}</h2>
      </div>
      <div className="profile-fields">
        <label>
          Profile
          <select value={selectedProfileId} onChange={(e) => selectProfile(e.target.value)}>
            {profiles.map((p) => (
              <option key={p.id} value={p.id}>
                {p.name}
              </option>
            ))}
          </select>
        </label>
        <label>
          Codec
          <select value={selected?.codec ?? 'hevc'} onChange={(e) => updateSelectedProfile({ codec: e.target.value as 'h264' | 'hevc' | 'av1' })}>
            <option value="h264">H.264</option>
            <option value="hevc">HEVC</option>
            <option value="av1">AV1</option>
          </select>
        </label>
        <label>
          Rate control
          <select
            value={selected?.rate_control ?? 'qvbr'}
            onChange={(e) => updateSelectedProfile({ rate_control: e.target.value as 'qvbr' | 'cqp' | 'cbr' | 'vbr' })}
          >
            <option value="qvbr">QVBR</option>
            <option value="cqp">CQP</option>
            <option value="cbr">CBR</option>
            <option value="vbr">VBR</option>
          </select>
        </label>
        <label>
          Rate value
          <input
            type="number"
            value={selected?.rate_value ?? 28}
            onChange={(e) => updateSelectedProfile({ rate_value: Number(e.target.value) })}
          />
        </label>
        <label>
          Custom options
          <textarea
            rows={3}
            value={selected?.custom_options ?? ''}
            onChange={(e) => updateSelectedProfile({ custom_options: e.target.value })}
          />
        </label>
      </div>
    </section>
  );
}
