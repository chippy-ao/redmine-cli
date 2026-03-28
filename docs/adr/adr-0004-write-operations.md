# ADR-0004: Write 操作の追加方式

## Status: Accepted (2026-03-28)

## Context

redmine-cli は当初 GET のみ対応の読み取り専用 CLI として設計した。
Claude Code Skills からチケット作成・リレーション操作を行うニーズが発生し、
POST/DELETE 操作の追加が必要になった。

HTTP クライアントの拡張方式として以下を検討した：

| 方法 | メリット | デメリット |
|------|----------|------------|
| doGet を doRequest にリファクタし共通基盤化 | ヘッダ設定・エラーハンドリングを共通化、既存テスト影響なし | リファクタの手間（1回のみ） |
| Post/Delete を独立実装 | 既存コードに触れない | コード重複、エラーハンドリングの二重管理 |
| 汎用 Do メソッド1つで全対応 | メソッド数が少ない | 型安全性が低い、呼び出し側が複雑 |

## Decision

doGet を doRequest にリファクタし、共通基盤の上に Post/Delete を追加する。

- `doRequest(*http.Request, any) error` が共通処理（ヘッダ設定、HTTP 実行、レスポンス処理）を担当
- `Get`, `GetRawQuery`, `Post`, `Delete` が各自リクエストを構築して `doRequest` に委譲
- 成功ステータスコードを 200/201/204 に拡張
- 422 (Unprocessable Entity) で Redmine の `{"errors": [...]}` をパース

コマンド構造はフラットを維持（`create-issue`, `add-relation`, `delete-relation`）。

## Consequences

- PR1 で1回だけリファクタ。以降の HTTP メソッド追加は1メソッドの追加で済む
- 既存の Get/GetRawQuery テストはリファクタ後もそのまま通る（doGet は private メソッド）
- コマンド数は増えるが、フラットのままで一貫性を保てる
