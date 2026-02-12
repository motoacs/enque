import { describe, expect, it } from 'vitest';
import { t } from './i18n';
import { useAppStore } from '../stores/appStore';

describe('i18n', () => {
  it('returns english string when language is en', () => {
    useAppStore.setState((state) => ({ config: { ...state.config, language: 'en' } }));
    expect(t('app.title')).toBe('Enque');
    expect(t('queue.title')).toBe('Queue');
  });
});
