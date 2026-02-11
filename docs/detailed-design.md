# Enque 詳細設計書（SSOT）

document version: v1.0

## 0. 文書情報

| 項目 | 値 |
| --- | --- |
| 文書ID | ENQ-DS-001 |
| 文書名 | Enque 詳細設計書（Single Source of Truth） |
| 対象計画書 | `docs/project-plan.md`（Enque プロジェクト計画書 v4.1） |
| 対象バージョン | Enque v1.0.0（NVEncC実装 + マルチエンコーダ拡張基盤） |
| 最終更新 | 2026-02-11 |
| 対象環境 | Windows 11 (x64), NVEncC 8.x 以降（QSVEncC/ffmpeg拡張を考慮） |

## 1. SSOT運用ルール

本書は Enque の実装・テスト・レビュー時の唯一の設計基準とする。

1. 要件解釈が曖昧な場合は、本書の定義を優先する。
2. `docs/project-plan.md` は「要求意図」、本書は「実装契約」の位置づけとする。
3. 本書変更時は、変更理由・影響範囲・移行方針を同時に更新する。
4. 実装が本書と乖離する場合、実装を修正するか、本書を先に改訂する。
5. すべての機能は要件ID（`US-*`, `F-*`）へ追跡可能であること。

## 2. 目的・範囲・非対象

### 2.1 目的

エンコーダCLI上級ユーザー向けに、以下を安定提供する。

- 複数動画の一括エンコード
- GUI設定 + カスタムオプションの両立
- 進捗監視、停止/中止/キャンセル制御
- ファイル日時復元を含むメタデータ保持
- 実行再現性を担保するジョブ記録（`job.json` + stderrログ）
- 将来のQSVEncC/ffmpeg対応でコア再利用できる拡張基盤

### 2.2 対象範囲

- Windows デスクトップアプリ（Wails v2）
- Go バックエンド（エンコーダ実行管理）
- React + TypeScript フロントエンド（編集UI / 実行UI）
- JSON 永続化（`%APPDATA%/Enque`）

### 2.3 非対象（v1）

- 動画プレビュー
- カット/結合/分割など編集機能
- ffprobe ベースの入力詳細解析
- コーデックとコンテナの互換性検証
- ジョブ個別プロファイル割り当て
- プロセスサスペンドによる一時停止
- QSVEncC/ffmpeg の実行機能そのもの（v1では未実装、拡張ポイントのみ定義）

## 3. システム構成

### 3.1 論理アーキテクチャ

```text
React(UI) + Zustand(store)
  ├─ Wails Binding (Request/Response)
  └─ Wails Events  (Backend -> Frontend)

Go Backend
  ├─ app.go                 バインディング公開面
  ├─ backend/queue          セッションとワーカープール
  ├─ backend/encoder        adapter解決・コマンド生成・実行・進捗パース
  ├─ backend/profile        プロファイルCRUD/マイグレーション
  ├─ backend/config         アプリ設定CRUD/マイグレーション
  ├─ backend/detector       NVEncC/QSVEncC/ffmpeg 検出とGPU情報
  ├─ backend/metadata       Win32 FileTime 復元
  └─ backend/logging        job.json と stderr 永続化

External Process
  └─ selected encoder binary x N (v1: NVEncC64.exe)
```

### 3.2 Source of Truth 切替

- エンコード前: Frontend (`editStore`, `profileStore`) が SoT
- エンコード開始時: スナップショットを `StartEncode` に送信
- エンコード中: Backend (`queue.Manager`) が SoT
- 完了後: Backend から最終結果を受け取り Frontend 表示へ反映

この切替により、実行中に編集状態が変化しても実行結果がぶれない。

### 3.3 並列化モデル

- レイヤ1（プロセス並列）: `max_concurrent_jobs`（1..8）
- レイヤ2（ジョブ内並列）: adapter定義（v1 `nvencc`: `--split-enc`, `--parallel`）

デフォルトは `max_concurrent_jobs=1` + `split_enc=auto`（`nvencc`）。
出力パス確定はジョブ開始直前に mutex 下で行い、並列時の競合を防ぐ。

## 4. ディレクトリ/モジュール設計

```text
backend/
  app/
    app.go
  queue/
    manager.go
    session.go
    output_resolver.go
  encoder/
    registry.go
    adapter.go
    nvencc/
      command_builder.go
      progress_parser.go
      capabilities.go
    qsvenc/
      command_builder.go
      progress_parser.go
      capabilities.go
    ffmpeg/
      command_builder.go
      progress_parser.go
      capabilities.go
    process_runner.go
    timeout_guard.go
  profile/
    model.go
    manager.go
    migration.go
  config/
    model.go
    manager.go
    migration.go
  detector/
    nvencc.go
    qsvenc.go
    ffmpeg.go
    tools.go
  metadata/
    file_time_windows.go
  logging/
    job_record.go
    stderr_writer.go
    app_logger.go
frontend/
  src/
    stores/
      editStore.ts
      profileStore.ts
      encodeStore.ts
      appStore.ts
    features/
      queue/
      profile/
      encode/
      settings/
    lib/
      api.ts
      events.ts
      i18n.ts
```

## 5. データ契約（厳密仕様）

## 5.1 列挙型

| 型 | 値 |
| --- | --- |
| `Codec` | `h264`, `hevc`, `av1` |
| `EncoderType` | `nvencc`, `qsvenc`, `ffmpeg` |
| `RateControl` | `qvbr`, `cqp`, `cbr`, `vbr` |
| `Preset` | `P1`..`P7` |
| `Multipass` | `none`, `quarter`, `full` |
| `SplitEnc` | `off`, `auto`, `auto_forced`, `forced_2`, `forced_3`, `forced_4` |
| `ParallelMode` | `off`, `auto`, `2`, `3` |
| `Decoder` | `avhw`, `avsw` |
| `AudioMode` | `copy`, `aac`, `opus` |
| `OnError` | `skip`, `stop` |
| `OverwriteMode` | `ask`, `auto_rename` |
| `PostAction` | `none`, `shutdown`, `sleep`, `custom` |
| `Language` | `ja`, `en` |
| `JobStatus` | `pending`, `running`, `completed`, `failed`, `cancelled`, `timeout`, `skipped` |

## 5.2 Profile スキーマ

計画書 6.8 を基準に、v1で以下を確定値とする。

```json
{
  "id": "uuid",
  "version": 2,
  "name": "string(1..80)",
  "is_preset": false,
  "encoder_type": "nvencc",
  "encoder_options": {},
  "codec": "hevc",
  "rate_control": "qvbr",
  "rate_value": 28,
  "preset": "P4",
  "output_depth": 10,
  "multipass": "none",
  "output_res": "",
  "bframes": null,
  "ref": null,
  "lookahead": null,
  "gop_len": null,
  "aq": true,
  "aq_temporal": true,
  "split_enc": "auto",
  "parallel": "off",
  "decoder": "avhw",
  "device": "auto",
  "audio_mode": "copy",
  "audio_bitrate": 256,
  "colormatrix": "auto",
  "transfer": "auto",
  "colorprim": "auto",
  "colorrange": "auto",
  "max_cll": "",
  "master_display": "",
  "dhdr10_info": "off",
  "dolby_vision_rpu": "off",
  "metadata_copy": true,
  "video_metadata_copy": true,
  "audio_metadata_copy": true,
  "chapter_copy": true,
  "sub_copy": true,
  "data_copy": true,
  "attachment_copy": true,
  "restore_file_time": false,
  "custom_options": ""
}
```

### 5.2.1 バリデーション

| 項目 | ルール |
| --- | --- |
| `name` | 必須、前後空白trim後 1..80 |
| `encoder_type` | `nvencc` / `qsvenc` / `ffmpeg` |
| `encoder_options` | JSON object（最大 64KB） |
| `rate_value` | `> 0` |
| `output_depth` | `8` or `10` |
| `bframes` | `null` or `0..7` |
| `lookahead` | `null` or `0..32` |
| `audio_bitrate` | `32..1024` |
| `custom_options` | UTF-8、最大 4096 文字 |
| `device` | `auto` or `0..15` |

バリデーション失敗時は保存不可とし、UI上に該当フィールドエラーを表示する。

互換ルール:

- `encoder_type=nvencc` の場合、本章で定義するNVEncC向けフィールド（`codec` 〜 `custom_options`）を有効化する
- `encoder_type!=nvencc` の場合、NVEncC向けフィールドは保存はするが実行時には無視し、`encoder_options` をadapterへ渡す

## 5.3 AppConfig スキーマ

計画書 6.9 を基準に、v1で以下を確定値とする。

```json
{
  "version": 1,
  "nvencc_path": "",
  "qsvenc_path": "",
  "ffmpeg_path": "",
  "ffprobe_path": "",
  "max_concurrent_jobs": 1,
  "on_error": "skip",
  "decoder_fallback": false,
  "keep_failed_temp": false,
  "no_output_timeout_sec": 600,
  "no_progress_timeout_sec": 300,
  "post_complete_action": "none",
  "post_complete_command": "",
  "output_folder_mode": "same_as_input",
  "output_folder_path": "",
  "output_name_template": "{name}_encoded.{ext}",
  "output_container": "mkv",
  "overwrite_mode": "ask",
  "language": "ja",
  "default_profile_id": ""
}
```

### 5.3.1 バリデーション

| 項目 | ルール |
| --- | --- |
| `max_concurrent_jobs` | `1..8` |
| `no_output_timeout_sec` | `30..86400` |
| `no_progress_timeout_sec` | `30..86400` |
| `output_name_template` | 必須、1..255、`{name}` または固定文字列を含むこと |
| `output_folder_path` | `output_folder_mode=specified` のとき必須 |
| `post_complete_command` | `post_complete_action=custom` のとき必須 |

## 5.4 実行時データ

### 5.4.1 QueueJob

```json
{
  "job_id": "uuid",
  "input_path": "C:\\path\\input.mp4",
  "input_size_bytes": 123456789,
  "status": "pending",
  "progress": {
    "percent": 0,
    "fps": null,
    "bitrate_kbps": null,
    "eta_sec": null
  },
  "started_at": "",
  "finished_at": "",
  "worker_id": null,
  "exit_code": null,
  "error_message": ""
}
```

### 5.4.2 EncodeSession

```json
{
  "session_id": "uuid",
  "state": "running",
  "started_at": "ISO8601",
  "finished_at": "",
  "total_jobs": 10,
  "completed_jobs": 2,
  "running_jobs": 1,
  "failed_jobs": 0,
  "cancelled_jobs": 0,
  "timeout_jobs": 0,
  "skipped_jobs": 0,
  "stop_requested": false,
  "abort_requested": false
}
```

## 5.5 JobRecord スキーマ

`docs/project-plan.md` 6.11 に準拠する。`schema_version=1` で固定。

## 6. API契約（Wails Binding）

## 6.1 一覧

| メソッド | 用途 |
| --- | --- |
| `Bootstrap()` | 起動時に設定・プロファイル・ツール検出結果を取得 |
| `SaveAppConfig(config)` | 設定保存 |
| `ListProfiles()` | プロファイル一覧取得 |
| `UpsertProfile(profile)` | プロファイル作成/更新 |
| `DeleteProfile(profileID)` | プロファイル削除 |
| `DuplicateProfile(profileID, newName)` | 複製 |
| `SetDefaultProfile(profileID)` | デフォルト設定 |
| `GetGPUInfo()` | `--check-device`, `--check-features` 結果取得 |
| `DetectExternalTools()` | NVEncC/QSVEncC/ffmpeg/ffprobe 検出 |
| `StartEncode(request)` | セッション開始 |
| `RequestGracefulStop(sessionID)` | 停止（次ジョブ抑止） |
| `RequestAbort(sessionID)` | 中止（実行中含め強制終了） |
| `CancelJob(sessionID, jobID)` | 実行中ジョブの個別中止 |
| `ResolveOverwrite(sessionID, jobID, decision)` | `overwrite_mode=ask` 応答 |
| `ListTempArtifacts()` | 残存 tmp 候補一覧取得 |
| `CleanupTempArtifacts(paths)` | 指定 tmp の削除 |

## 6.2 StartEncode 入力契約

```json
{
  "jobs": [{"job_id": "uuid", "input_path": "C:\\in.mp4"}],
  "profile": {"...": "Profile object"},
  "app_config_snapshot": {"...": "AppConfig object"},
  "command_preview": "optional string shown in UI"
}
```

ルール:

1. `jobs` は 1 件以上必須。
2. 実行中セッションがある場合は `E_SESSION_RUNNING` を返す。
3. `profile` と `app_config_snapshot` は保存済み値ではなく、開始時点スナップショットを使用する。
4. `profile.encoder_type` に対応するadapterが未登録の場合は `E_ENCODER_NOT_IMPLEMENTED` を返す。

## 6.3 エラーコード

| コード | 意味 | UI動作 |
| --- | --- | --- |
| `E_VALIDATION` | 入力バリデーション違反 | フィールドエラー表示 |
| `E_TOOL_NOT_FOUND` | 選択エンコーダ未検出 | 設定画面へ誘導 |
| `E_TOOL_VERSION_UNSUPPORTED` | 選択エンコーダの対象バージョン外 | 警告表示し開始不可 |
| `E_ENCODER_NOT_IMPLEMENTED` | 対象adapter未実装 | 通知 + encoder変更導線 |
| `E_SESSION_RUNNING` | 既に実行中セッションあり | 二重開始防止 |
| `E_IO` | JSON保存/ログ保存失敗 | 通知 + 再試行導線 |
| `E_INTERNAL` | 予期しない内部エラー | 通知 + ログ参照導線 |

## 7. Event契約（Go -> Frontend）

イベント名は `enque:*` の名前空間で統一する。

| イベント | payload 概要 |
| --- | --- |
| `enque:session_started` | `session_id`, `total_jobs`, `started_at`, `encoder_type` |
| `enque:job_started` | `session_id`, `job_id`, `worker_id`, `input_path`, `temp_output_path`, `encoder_type` |
| `enque:job_progress` | `session_id`, `job_id`, `percent`, `fps`, `bitrate_kbps`, `eta_sec`, `raw_line` |
| `enque:job_log` | `session_id`, `job_id`, `line`, `ts` |
| `enque:job_needs_overwrite` | `session_id`, `job_id`, `final_output_path` |
| `enque:job_finished` | `session_id`, `job_id`, `status`, `exit_code`, `error_message`, `final_output_path` |
| `enque:session_state` | 集計値（completed/failed/...） |
| `enque:session_finished` | セッション最終結果 |
| `enque:warning` | 非致命警告（パース失敗など） |
| `enque:error` | 致命エラー |

### 7.1 `job_progress` スロットリング

- 最低 500ms 間隔で送信
- 最終行（ジョブ終了直前）は即時送信
- パース不可行は `percent=null` で送信可（ログ表示を優先）

## 8. フロントエンド設計

## 8.1 Store責務

| Store | 責務 |
| --- | --- |
| `editStore` | キュー編集、選択中プロファイル、開始前UI状態 |
| `profileStore` | プロファイルCRUD、保存状態、バリデーション |
| `encodeStore` | セッション状態、ジョブ進捗、ログリングバッファ |
| `appStore` | AppConfig、ツール検出結果、GPU情報、言語 |

## 8.2 画面/コンポーネント

1. キューパネル
2. プロファイル編集パネル
3. 出力設定パネル
4. エンコーダ選択 + エンコーダ固有設定パネル
5. コマンドプレビュー + コピー
6. 実行コントロール（開始/停止/中止）
7. 実行モニタ（ジョブ進捗、全体進捗、ログ）
8. 設定ダイアログ（外部ツール、同時実行数、エラー時挙動）
9. GPU情報ダイアログ（v1: NVEncC）

## 8.3 UI状態遷移

| 現在状態 | イベント | 次状態 |
| --- | --- | --- |
| `idle` | 開始成功 | `running` |
| `running` | 停止要求 | `stopping` |
| `running/stopping` | 全ジョブ終了 | `completed` |
| `running/stopping` | 中止要求 | `aborting` |
| `aborting` | 実行中ジョブ終了 | `aborted` |

ルール:

- `running` 中はキュー編集をロック
- 停止は次ジョブ開始を抑止
- 中止は実行中ジョブを強制終了

## 8.4 コマンドプレビュー

- 生成元は Backend の `encoder adapter` と同一ロジック
- プレビューは表示専用（ベストエフォート）
- クリップボードコピー機能を提供

## 9. バックエンド実装設計

## 9.1 Queue Manager

### 9.1.1 主要構造体

- `Manager`
- `Session`
- `Worker`
- `JobRuntime`

### 9.1.2 実行アルゴリズム（擬似コード）

```go
StartEncode(request):
  validate(request)
  adapter = registry.resolve(request.profile.encoder_type)
  session = newSession(snapshot)
  emit(session_started)
  spawn workers(max_concurrent_jobs)
  enqueue all pending jobs
  wait until all jobs done or abort
  run post_complete_action if eligible
  emit(session_finished)
```

Worker:

```go
for job in queue:
  if session.stopRequested { mark skipped; continue }
  runJob(job)
  if failed && on_error == stop { session.stopRequested = true }
```

## 9.2 出力パス確定

### 9.2.1 テンプレート変数

| 変数 | 意味 |
| --- | --- |
| `{name}` | 入力ファイルの拡張子除去名 |
| `{ext}` | 出力コンテナ拡張子（例: `mkv`） |

v1では上記2種類のみを正式サポートする。未知変数は文字列として残す。

### 9.2.2 確定手順

1. 出力ディレクトリ決定（入力同一 or 指定）
2. テンプレート適用して基底名を生成
3. `overwrite_mode` 判定
4. `auto_rename` の場合、mutex 下で `_001`, `_002`... を採番
5. temp 出力を `{name}.{short_id}.tmp.{ext}` として確定
6. 成功時に temp -> final へ rename（同一ボリューム前提で原子的）

### 9.2.3 `overwrite_mode=ask`

- 既存ファイル衝突時に `enque:job_needs_overwrite` を発火
- Frontend は `overwrite` / `skip` / `abort` を選択
- 応答は `ResolveOverwrite` で返す
- 応答待ちタイムアウトは 10 分（超過時 `skip`）

## 9.3 Command Builder（adapter）

## 9.3.1 引数生成順序（固定）

1. adapter前置オプション
2. 入力
3. GUI層オプション
4. カスタムオプション（後勝ち）
5. 出力

`nvencc` adapter では従来順序（`--avhw/--avsw` -> `-i` -> `-c` -> GUI -> custom -> `-o`）を固定契約として維持する。順序変更は互換性破壊として扱う。

## 9.3.2 GUI項目 -> NVEncC オプション（`nvencc` adapter）

| GUI項目 | 条件 | 出力引数 |
| --- | --- | --- |
| `rate_control=qvbr` | 常時 | `--qvbr <rate_value>` |
| `rate_control=cqp` | 常時 | `--cqp <rate_value>` |
| `rate_control=cbr` | 常時 | `--cbr <rate_value>` |
| `rate_control=vbr` | 常時 | `--vbr <rate_value>` |
| `preset` | 常時 | `--preset <P1..P7>` |
| `output_depth` | 常時 | `--output-depth <8|10>` |
| `multipass!=none` | 条件 | `--multipass <quarter|full>` |
| `output_res!=empty` | 条件 | `--output-res <WxH>` |
| `bframes!=null` | 条件 | `--bframes <n>` |
| `ref!=null` | 条件 | `--ref <n>` |
| `lookahead!=null` | 条件 | `--lookahead <n>` |
| `gop_len!=null` | 条件 | `--gop-len <n>` |
| `aq=true` | 条件 | `--aq` |
| `aq_temporal=true` | 条件 | `--aq-temporal` |
| `split_enc!=off` | 条件 | `--split-enc <mode>` |
| `parallel!=off` | 条件 | `--parallel <mode>` |
| `device!=auto` | 条件 | `--device <id>` |
| `audio_mode=copy` | 常時 | `--audio-copy` |
| `audio_mode=aac` | 常時 | `--audio-codec aac --audio-bitrate <n>` |
| `audio_mode=opus` | 常時 | `--audio-codec opus --audio-bitrate <n>` |
| `colormatrix!=auto` | 条件 | `--colormatrix <value>` |
| `transfer!=auto` | 条件 | `--transfer <value>` |
| `colorprim!=auto` | 条件 | `--colorprim <value>` |
| `colorrange!=auto` | 条件 | `--colorrange <value>` |
| `max_cll!=empty` | 条件 | `--max-cll <value>` |
| `master_display!=empty` | 条件 | `--master-display <value>` |
| `dhdr10_info=copy` | 条件 | `--dhdr10-info copy` |
| `dolby_vision_rpu=copy` | 条件 | `--dolby-vision-rpu copy` |
| `metadata_copy=true` | 条件 | `--metadata copy` |
| `video_metadata_copy=true` | 条件 | `--video-metadata copy` |
| `audio_metadata_copy=true` | 条件 | `--audio-metadata copy` |
| `chapter_copy=true` | 条件 | `--chapter-copy` |
| `sub_copy=true` | 条件 | `--sub-copy` |
| `data_copy=true` | 条件 | `--data-copy` |
| `attachment_copy=true` | 条件 | `--attachment-copy` |

## 9.3.3 カスタムオプションの字句解析

- `custom_options` は quote 対応トークナイザで分割する（adapter共通の既定実装）。
- 対応: `"..."`, `'...'`, `\"`, `\'`
- 不正クォート時は `E_VALIDATION` で開始不可。
- 分割後トークンを引数末尾へ追加する（後勝ち）。adapterが独自パーサを要求する場合は上書き可能。

## 9.4 Process Runner

## 9.4.1 起動

- `exec.CommandContext` で起動
- stderr パイプを逐次読取
- 起動成功後に Job Object を割り当てる

## 9.4.2 Job Object

1. `CreateJobObject`
2. `SetInformationJobObject(JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE)`
3. `AssignProcessToJobObject`
4. キャンセル時 `TerminateJobObject`
5. 失敗時 `taskkill /F /T /PID` へフォールバック

`used_job_object` を `job.json` に記録する。

## 9.4.3 ハング検出

- `lastLineAt` と `lastProgressAt` を保持
- 1秒 ticker で監視
- `now-lastLineAt > no_output_timeout_sec` で timeout
- `progress有効 && now-lastProgressAt > no_progress_timeout_sec` で timeout

## 9.4.4 デコーダフォールバック

再試行条件:

1. `decoder_fallback=true`
2. 初回実行の実効デコーダが `avhw`
3. 終了コードが非0
4. 再試行未実施

条件を満たす場合、`--avsw` で1回だけ再実行し `retry_applied=true`, `retry_detail="nvencc: avhw->avsw"` を記録する。
このフォールバックは `nvencc` adapter のみ対象とする。

## 9.5 進捗パーサ

- parserはadapter単位で実装する。v1で本契約を満たすのは `nvencc` parser。
- 区切り: `\r` と `\n`
- 正規表現は複数パターン許容（コーデック差異に対応）
- 抽出項目: `percent`, `fps`, `bitrate_kbps`, `eta_sec`
- パース失敗時は `percent=null` で処理継続

## 9.6 ファイル日時復元

- 対象: `CreationTime`, `LastWriteTime`
- 非対象: `LastAccessTime`
- 実行タイミング: temp -> final rename 成功後
- 失敗時: ジョブは成功扱い、警告をログ出力

## 9.7 完了後アクション

- 実行条件: セッションが `aborted` でないこと
- 種別:
  - `none`: 何もしない
  - `shutdown`: Windowsシャットダウンコマンド
  - `sleep`: Windowsスリープコマンド
  - `custom`: `post_complete_command` 実行
- `custom` はコマンド文字列をそのまま実行し、終了コードをアプリログに残す

## 10. 永続化設計

## 10.1 保存先

`%APPDATA%/Enque/` 配下に以下を配置する。

```text
%APPDATA%/Enque/
  config.json
  profiles.json
  logs/
    {job_id}.json
    {job_id}.stderr.log
  runtime/
    temp_index.json
```

## 10.2 JSON 書き込み

すべて原子的に保存する。

1. `*.tmp` に書き込み
2. flush + close
3. rename で置換

## 10.3 マイグレーション

- `config.version`, `profile.version` で判定
- 旧版読み込み時は `migration.go` で最新へ変換
- 変換不能時はバックアップ保存後、デフォルト生成 + 警告表示
- `profile` の v1 -> v2 変換では `encoder_type="nvencc"` と `encoder_options={}` を補完する

## 10.4 残存 tmp 検出

- 実行中に生成した temp パスを `runtime/temp_index.json` に追記
- 正常終了/削除時に index から除去
- 起動時に index を検査し、存在する temp を「削除候補」として提示

## 11. 外部ツール検出設計

## 11.1 NVEncC

検索順:

1. `config.nvencc_path`
2. アプリ実行ファイル同一ディレクトリ（`NVEncC64.exe`, `NVEncC.exe`）
3. `PATH`

検出後:

- `--version` またはヘッダ行でバージョン取得
- 8.x 未満は `E_TOOL_VERSION_UNSUPPORTED`

## 11.2 QSVEncC

検索順:

1. `config.qsvenc_path`
2. アプリ実行ファイル同一ディレクトリ（`QSVEncC64.exe`, `QSVEncC.exe`）
3. `PATH`

v1では未検出でも実行可能とする（警告のみ）。将来のadapter有効化時に必須化する。

## 11.3 ffmpeg / ffprobe

同様に検索するが、未検出でも実行可能とする（警告のみ）。

## 12. i18n設計

- 初期対応言語: `ja`, `en`
- 翻訳キーは機能単位で管理（例: `encode.start`, `queue.clear`）
- エラーメッセージは backend でコード化し frontend でローカライズ

## 13. ログ設計

## 13.1 UIログ（メモリ）

- ジョブ単位リングバッファ
- デフォルト保持行数: 2000 行/ジョブ
- 古い行から破棄

## 13.2 永続ログ（ファイル）

- stderr は全文保存（省略なし）
- `job.json` に再現用情報（argv, exit_code, retry有無, worker_id）を保存

## 13.3 アプリログ

- `app.log`（日次ローテーション、30日保持）
- 重大イベント（開始/停止/中止/フォールバック/タイムアウト）を記録

## 14. セキュリティ・安全性

1. `custom_options` と `post_complete_command` はユーザー責任の上級機能として扱う。
2. 文字列はログ保存時に制御文字をエスケープし、ログ破壊を防止する。
3. 削除操作（tmp cleanup）はユーザー選択対象のみに限定する。
4. 実行ファイルパスは存在確認 + 実行権限確認を行う。

## 15. テスト設計

## 15.1 ユニットテスト（必須）

| 対象 | 観点 |
| --- | --- |
| `encoder/registry.go` | `encoder_type` と adapter解決、未実装エラー |
| `nvencc/command_builder.go` | 引数順序、省略条件、後勝ち、Windowsパス |
| `nvencc/progress_parser.go` | 正常系/異常系、コーデック別、split/parallel時 |
| `output_resolver.go` | auto_rename採番、mutex排他、テンプレート |
| `profile/config migration` | 旧版JSON -> 最新版 |
| `timeout_guard.go` | no_output/no_progress 判定 |

## 15.2 結合テスト

1. `nvencc` モックで進捗イベント連携確認
2. 複数ワーカー実行時の状態整合性確認
3. 停止/中止/キャンセルの相互干渉確認
4. overwrite `ask` 応答待ちと再開確認
5. Job Object 失敗時フォールバック確認
6. `encoder_type` 不一致時に `E_ENCODER_NOT_IMPLEMENTED` を返すこと

## 15.3 手動テストマトリクス

- コーデック: H.264 / HEVC / AV1
- レート制御: QVBR / CQP / CBR / VBR
- split-enc: off / auto / forced_3
- parallel: off / auto / 2 / 3
- 同時実行数: 1 / 2 / 3 / 4
- エラー系: 不正入力、未対応オプション、NVEncC未検出、タイムアウト、未実装encoder_type

## 16. 受け入れ基準（DoD）

1. `F-01`〜`F-14` を満たす実装とテスト結果が揃っている。
2. 進捗が取得不能でもジョブ完走し、ログに生出力が残る。
3. 並列実行時に最終出力名衝突が発生しない。
4. `job.json` と `stderr.log` が全ジョブで生成される。
5. `nvencc` 選択時、NVEncC 8.x 未満で警告表示され、開始を拒否する。
6. ファイル日時復元ON時、CreationTime/LastWriteTimeが入力と一致する。
7. `encoder_type` が `qsvenc` / `ffmpeg` でadapter未実装の場合、`E_ENCODER_NOT_IMPLEMENTED` で安全に拒否する。

## 17. 要件トレーサビリティ

| 要件ID | 実装章 |
| --- | --- |
| `US-01` `US-09` `US-10` `F-01` | 8, 9.1 |
| `US-02` `US-05` `US-06` `F-02` | 5.2, 6, 8 |
| `US-07` `F-04` | 9.3.3 |
| `US-08` `F-05` | 8.4, 9.3 |
| `F-03` | 9.3.2 |
| `F-06` | 9.2, 5.3 |
| `US-03` `F-07` | 6.2, 7, 9 |
| `US-11` `US-12` `F-08` | 8.3, 9.1, 9.4 |
| `US-04` `F-09` | 9.6 |
| `US-13` `F-10` | 9.7 |
| `US-15` `F-11` | 6.1, 11 |
| `US-16` `F-12` | 11 |
| `US-14` | 3.3, 9.1 |
| `F-13` | 5.2（プリセット定義）, 6（複製運用） |
| `US-17` `F-14` | 3.1, 4, 5.2, 9.3 |

## 18. 同梱プリセット定義（固定値）

v1 同梱プリセットは以下とする（`is_preset=true`, `encoder_type="nvencc"`）。

1. `HEVC Quality`
2. `AV1 Fast`
3. `Camera Archive`
4. `H.264 Compatible`

値は計画書記載内容に従う。プリセット本体は直接編集不可、複製して編集する運用とする。

## 19. 未解決事項（v1実装前に確定）

1. `overwrite_mode=ask` のUI（ジョブ単位ダイアログ or 一括ダイアログ）
2. `post_complete_action=sleep` の実行コマンド最終確定
3. `nvencc` 進捗正規表現の最終版（実機ログ収集後）
4. `qsvenc` / `ffmpeg` adapter のオプションスキーマ確定
5. `app.log` ローテーション実装方式（自前 or ライブラリ）

上記は実装着手時の最初の技術決定事項とし、確定後に本書を更新する。
