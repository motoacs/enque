import { useAppStore } from '../stores/appStore';

const messages = {
  ja: {
    'app.title': 'Enque',
    'queue.title': 'ジョブキュー',
    'profile.title': 'プロファイル',
    'settings.title': '設定',
    'encode.title': '実行モニタ',
    'button.start': '開始',
    'button.stop': '停止',
    'button.abort': '中止',
    'button.copy': 'コピー',
    'preview.title': 'コマンドプレビュー'
  },
  en: {
    'app.title': 'Enque',
    'queue.title': 'Queue',
    'profile.title': 'Profile',
    'settings.title': 'Settings',
    'encode.title': 'Encode Monitor',
    'button.start': 'Start',
    'button.stop': 'Stop',
    'button.abort': 'Abort',
    'button.copy': 'Copy',
    'preview.title': 'Command Preview'
  }
} as const;

export function t(key: keyof (typeof messages)['ja']): string {
  const lang = useAppStore.getState().config.language;
  return messages[lang][key] ?? messages.ja[key];
}
