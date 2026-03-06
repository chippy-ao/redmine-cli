# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Redmine REST API を操作する Go 製 CLI ツール `redmine-cli`。
Claude Code Skills からの利用を主なユースケースとし、全コマンドが JSON を stdout に出力する。
Redmine 5.1.0+ が必要（`any_searchable` フィルタ使用）。

## Build & Test Commands

```bash
go build -o redmine-cli .    # ビルド
go test ./...                 # 全テスト実行
go test ./internal/query/...  # 単一パッケージのテスト
go vet ./...                  # 静的解析
```

Go のバージョンは mise で管理（`mise.toml`）。

## Architecture

3層構造: `cmd/` → `internal/` → Redmine REST API

- **`cmd/`** — cobra コマンド定義。各コマンドは `loadClientFromProfile()` でクライアント取得、`outputJSON()` で結果出力する共通パターン。
- **`internal/config/`** — YAML プロファイル管理。設定ファイルは `~/.config/redmine-cli/config.yaml`。パーミッションはディレクトリ 0700、ファイル 0600。
- **`internal/client/`** — Redmine API HTTP クライアント。`Get()` は通常パラメータ、`GetRawQuery()` はエンコード済みクエリ文字列をそのまま渡す。
- **`internal/query/`** — 検索クエリ構築。keyword 検索時は Redmine のフィルタ構文 (`f[]/op[]/v[]`) に変換する。`status_id` の特殊値 (open/closed/*) はフィルタ形式に含めず通常パラメータとして渡す。

## Key Design Decisions

- `--keyword` 指定時と非指定時で API 呼び出し方法が分岐する（`cmd/search.go`）。keyword あり → `GetRawQuery` でフィルタ構文、なし → `Get` で通常パラメータ。
- `--profile` グローバルフラグで複数 Redmine 環境を切り替え可能。
- 出力は v1.0 では JSON のみ。エラーは stderr + 終了コード 1。
- 設計判断は `docs/adr/` に ADR として記録。アーキテクチャ全体は `docs/architecture.md` を参照。

## CI/CD

- `.github/workflows/test.yaml` — push/PR で `go vet` + `go test`
- `.github/workflows/release.yaml` — `v*` タグで goreleaser 実行（linux/darwin/windows × amd64/arm64）
- Homebrew tap (`chippy-ao/homebrew-tap`) は goreleaser が自動更新
