# 次の実装スコープと計画（OAuth AS PoC）

## スコープ（次フェーズ）
- 認可UIの最小実装: ダミーログイン＋同意画面（同意スコープ表示、許可/拒否）。
- リダイレクトURI検証の強化: 127.0.0.1/::1 の可変ポート許容＋パス固定、その他は完全一致のみ。
- リソース宛先の厳格化: `resource` 未指定時のデフォルト抑制、`aud`/`resource` の検証方針を明文化。
- サンプル保護リソースAPI: `GET /resource/echo`（Bearer必須、`WWW-Authenticate` 応答をRFC 6750で検証可能に）。
- クライアント支援スクリプト（最小）: PKCE＋state生成、認可URL生成、手動貼付URLの解析、トークン交換の一括実行（CLI）。

## タスク一覧
- UI: `web/templates/login.html`, `web/templates/consent.html` とハンドラ追加（セッションはメモリ）。
- バリデーション: `internal/server/oauth/validator.go` の redirect 検証ロジック強化、テストデータ追加。
- リソースAPI: `cmd/as-poc` に `GET /resource/echo` を追加、`Authorization` 検証と `WWW-Authenticate` 実装。
- CLIツール: `scripts/auth-cli`（Go/シェルいずれか）で以下を提供：
  - `new`（PKCE＋state生成、JSON出力）
  - `url`（認可URL生成）
  - `complete`（フルURL貼付→code/state検証→トークン交換）
- ドキュメント: `poc-docs/oauth-as/README.md` に手動コピー方式の手順を追記。

## マイルストーン/順序
1) リダイレクト検証強化（安全性の底上げ）
2) サンプル保護リソースAPI追加（Bearer検証/挑戦の確認）
3) 認可UIの最小実装（ログイン/同意）
4) クライアント支援CLI（最小）
5) ドキュメント整備・動作例

## 受け入れ基準
- 認可→同意→302→手動貼付→トークン交換→保護リソース呼び出しの一連が手順書どおりに再現可能。
- 無効トークン/未提示時に `WWW-Authenticate: Bearer ...` を返却。
- 不正な `redirect_uri` は拒否、`state` 不一致/期限切れの再現と説明が可能。
