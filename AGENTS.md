# AGENTS.md

このファイルは、`/Users/yuta/Git/enque` 配下で作業するエージェント向けの実行ガイドです。  
実装時の判断を統一し、`docs/project-plan.md` と `docs/detailed-design.md` に沿った変更だけを行うことを目的とします。

## 1. 目的と優先順位

- 本リポジトリでは `docs/detailed-design.md` を SSOT（実装契約）として扱う。
- 要件意図は `docs/project-plan.md`、実装判断は `docs/detailed-design.md` を優先する。
- より深いディレクトリに別の `AGENTS.md` がある場合は、そちらを優先する。

## 2. スコープ

- 対象: Enque（Wails v2 / Go / React + TypeScript）アプリの設計・実装・テスト・ドキュメント更新。
- 非対象: 計画書にない機能追加（動画編集、プレビュー等）や、仕様未合意の大規模拡張。

## 3. 変更前チェック

実装開始前に以下を確認すること。

1. 変更対象機能の要件ID（`US-*`, `F-*`）を `docs/project-plan.md` で特定する。
2. 対応する実装章を `docs/detailed-design.md` の「要件トレーサビリティ」で確認する。
3. 既存仕様と矛盾する場合は、先にドキュメントを更新してからコードを変更する。

## 4. 実装ルール

- 仕様の新規追加・挙動変更を伴う変更では、実装と同時に `docs/detailed-design.md` を更新する。
- Wails バインディング API・イベント名・データスキーマは、設計書の契約を破らない。
- 並列実行、出力名確定、キャンセル/停止、タイムアウトは安全側（データ破壊回避）で実装する。
- ログは再現性重視: `job.json` と `stderr` の保存仕様を維持する。
- 不明点は推測実装せず、TODO/未解決事項として明示する。

## 5. コーディング規約

- Go:
  - `go fmt` 前提の整形を維持する。
  - エラーは握りつぶさず、呼び出し元で扱える形で返す。
  - 並行処理は `context.Context` とキャンセル伝播を必須とする。
- TypeScript/React:
  - 型安全を優先し、`any` の導入は最小限にする。
  - 状態管理は store の責務境界（`edit/profile/encode/app`）を越えて混在させない。
  - UI文言は i18n キー経由で扱う。

## 6. テスト方針

- 変更に対応するテストを追加または更新する（最低1件）。
- 優先対象:
  - `command_builder`
  - `progress_parser`
  - 出力パス解決（排他制御/連番）
  - 設定・プロファイルのマイグレーション
- バグ修正時は再発防止テストを必ず追加する。

## 7. ドキュメント更新ルール

- 機能仕様を変える PR では、次を同時更新する。
  - `docs/detailed-design.md`（必須）
  - `docs/project-plan.md`（要件意図が変わる場合のみ）
- 受け入れ条件（DoD）に影響がある場合、該当章を更新する。

## 8. 変更提案・PRメッセージの最小要件

変更説明には最低限以下を含める。

1. どの要件IDを満たす変更か
2. どの契約（API/Event/Schema）に影響するか
3. 互換性影響の有無
4. 追加・更新したテスト

## 9. 禁止事項

- 設計書契約と矛盾する API 追加/変更を、ドキュメント更新なしで行うこと
- 進捗パース失敗時にエンコード自体を停止させること（設計では継続）
- 実行中ジョブの中間生成物を最終出力として扱うこと
- 根拠のない最適化で並列制御の安全性を下げること

## 10. 現在のリポジトリ状態に関する注記

- 現在は主に仕様・設計ドキュメントが存在し、実装コードは未配置。
- そのため新規コード追加時は、先にディレクトリ構成・ビルド/テストコマンドを明示し、本ファイルにも追記すること。

## 11. 初期ディレクトリ構成・ビルド/テストコマンド

実装開始後は以下の構成を基準とする。

```text
backend/
  app/
  config/
  detector/
  encoder/
  logging/
  metadata/
  profile/
  queue/
frontend/
  src/
    features/
    lib/
    stores/
```

標準コマンド:

- Goテスト: `go test ./...`
- Go整形: `gofmt -w ./backend ./cmd`
- Frontend依存解決: `cd frontend && npm install`
- Frontendビルド: `cd frontend && npm run build`
- Frontendテスト: `cd frontend && npm run test`
- Wails開発起動: `wails dev`
- Wails本番ビルド: `wails build -platform windows/amd64`
