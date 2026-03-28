# Redmine CLI アーキテクチャ

## 概要

Redmine REST API を Bash から操作する Go 製 CLI ツール。
Claude Code Skills からの利用を主なユースケースとし、JSON 出力をデフォルトとする。

## 技術スタック

| 項目 | 選択 |
|------|------|
| 言語 | Go |
| CLI フレームワーク | [cobra](https://github.com/spf13/cobra) |
| 設定ファイル | YAML (`~/.config/redmine-cli/config.yaml`) |
| ビルド・リリース | goreleaser + GitHub Actions |
| テスト | Go 標準 `testing` パッケージ + `httptest` |

## プロジェクト構成

```
redmine-cli/
├── cmd/                     # cobra コマンド定義
│   ├── root.go              # ルートコマンド（--profile）
│   ├── search.go            # search コマンド
│   ├── get_issue.go         # get-issue コマンド
│   ├── projects.go          # projects コマンド
│   ├── trackers.go          # trackers コマンド
│   ├── statuses.go          # statuses コマンド
│   ├── categories.go        # categories コマンド
│   ├── versions.go          # versions コマンド
│   └── config.go            # config サブコマンド（add/list/remove/set-default）
├── internal/
│   ├── client/              # Redmine API クライアント
│   │   └── client.go
│   ├── config/              # 設定ファイル読み書き
│   │   └── config.go
│   └── query/               # 検索クエリ構築（keyword フィルタ等）
│       └── builder.go
├── main.go                  # エントリポイント
├── go.mod
├── .goreleaser.yaml
└── .github/workflows/
    ├── test.yaml
    └── release.yaml
```

`internal/` は Go の言語仕様により外部パッケージからインポート不可。内部実装の隠蔽に使用。

## コマンド体系

```
redmine-cli <command> [flags]

Commands:
  config       プロファイル管理（add / list / remove / set-default）
  search       チケット検索（キーワード + フィルタ）
  get-issue    チケット詳細取得
  projects     プロジェクト一覧
  trackers     トラッカー一覧
  statuses     ステータス一覧
  categories   カテゴリ一覧（要 --project）
  versions     バージョン一覧（要 --project）
  priorities      優先度一覧
  members         プロジェクトメンバー一覧（要 --project）
  create-issue    チケット作成（要 --project, --subject）
  add-relation    リレーション作成（要 --issue-id, --related-id, --type）
  delete-relation リレーション削除

Global Flags:
  --profile    使用するプロファイル名（省略時はデフォルト）
  --help       ヘルプ表示
```

### search フラグ

| フラグ | 型 | 説明 |
|-------|-----|------|
| `--keyword` | string | 全文テキスト検索（Redmine 5.1+） |
| `--project` | string | プロジェクト ID or 識別子 |
| `--status` | string | `open`, `closed`, `*`, または数値ID |
| `--assigned-to` | string | 担当者ID（`me` 可） |
| `--tracker-id` | int | トラッカーID |
| `--category-id` | int | カテゴリID |
| `--version-id` | int | 対象バージョンID |
| `--sort` | string | ソート列（例: `updated_on:desc`） |
| `--offset` | int | 取得開始位置（デフォルト: 0） |
| `--limit` | int | 取得件数（デフォルト: 25, 最大: 100） |

`--keyword` 指定時は Redmine のフィルタ構文 `f[]=any_searchable&op[any_searchable]=~&v[any_searchable][]=<keyword>` に変換。

### get-issue フラグ

| フラグ | 型 | 説明 |
|-------|-----|------|
| (引数) | int | チケットID（必須） |
| `--include` | string | 追加情報（カンマ区切り）: journals, children, relations, attachments, changesets, watchers, allowed_statuses |

### create-issue フラグ

| フラグ | 型 | 説明 |
|-------|-----|------|
| `--project` | string | プロジェクト ID or 識別子（必須） |
| `--subject` | string | チケット件名（必須） |
| `--tracker-id` | int | トラッカーID |
| `--status-id` | int | ステータスID |
| `--priority-id` | int | 優先度ID |
| `--description` | string | チケット説明 |
| `--category-id` | int | カテゴリID |
| `--version-id` | int | 対象バージョンID |
| `--assigned-to-id` | int | 担当者ID |
| `--parent-issue-id` | int | 親チケットID |
| `--estimated-hours` | float | 予定工数 |
| `--private` | bool | プライベートチケット |

未指定のオプションフラグはリクエストに含まれない（`cobra.Changed()` で制御）。

## 認証・設定

### プロファイル方式

```yaml
# ~/.config/redmine-cli/config.yaml
default_profile: work
profiles:
  work:
    url: https://redmine.company.com
    api_key: abc123
  oss:
    url: https://redmine.example.org
    api_key: def456
```

### config コマンド

```bash
redmine-cli config add <name> --url <url> --api-key <key>
redmine-cli config set-default <name>
redmine-cli config list
redmine-cli config remove <name>
```

## 出力仕様

- 成功: 終了コード 0 + stdout に JSON
- 失敗: 終了コード 1 + stderr にエラーメッセージ
- v1.0 は JSON のみ

## 配布方法

| OS | 方法 |
|----|------|
| macOS | `brew install chippy-ao/tap/redmine-cli` |
| macOS/Linux | `mise use -g chippy-ao/redmine-cli` |
| Windows | GitHub Releases からバイナリをダウンロード |
| Go 環境あり | `go install github.com/chippy-ao/redmine-cli@latest` |

リリースは goreleaser + GitHub Actions で自動化。タグ `v*` push で全OS向けバイナリをビルドし、GitHub Releases に公開、homebrew-tap の Formula を自動更新。

## エラーハンドリング

| エラー | stderr メッセージ | 終了コード |
|--------|-----------------|-----------|
| プロファイル未設定 | `プロファイルが設定されていません。redmine-cli config add で設定してください。` | 1 |
| 認証エラー (401) | `認証エラー: APIキーが無効です。` | 1 |
| 権限エラー (403) | `権限エラー: アクセス権がありません。` | 1 |
| 未検出 (404) | `チケット #xxx が見つかりません。` | 1 |
| バリデーションエラー (422) | `バリデーションエラー: Subject cannot be blank, ...` | 1 |

## 前提条件

- Redmine 5.1.0+（`any_searchable` フィルタ使用のため）
- ユーザーのランタイム依存なし（シングルバイナリ）
