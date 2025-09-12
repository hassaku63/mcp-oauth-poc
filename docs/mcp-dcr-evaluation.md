# MCP Server × Dynamic Client Registration — 評価

## ユーザーの問い
> MCP Server に dynamic client registration を組み込むことができると MCP のドキュメントに書かれている。
>
> ここで、OAuth client は各開発者のホストマシンにインストールされた MCP Client アプリケーションであるため、
> 不特定の多くの開発者向けに配布する MCP Server を (HTTP Transport を使って) 開発する場合は Dynamic client registration のユースケースが非常に適している、と予想する。
> ユーザーの問いは妥当か？評価して。

## 結論（妥当性）
- 妥当: MCP Server を HTTP Transport で不特定多数の開発者向けに配布し、各開発者のローカル MCP Client（ネイティブ/デスクトップ/CLI 等）を OAuth クライアントとして扱う設計において、Dynamic Client Registration（DCR; RFC 7591）はスケーラブルなクライアント・オンボーディング手段として適合する。
- 根拠: RFC 7591 は「事前の手作業/個別合意に依存する登録」を置き換えるために定義され（§1）、登録エンドポイントでメタデータを受け付け（§4）、発行結果・管理連携（§5, RFC 7592）まで含むため、大量・多様なクライアント配布の自動化に合う。

## 適合理由（ユースケース適性）
- セルフサービス登録: クライアントごとに `client_id` を自動払い出し（必要に応じて `client_secret` 以外の認証方式）し、配布/実行時のオンボーディングを自動化（RFC 7591 §1, §4）。
- ポリシー柔軟性: オープン登録 or 初期アクセス・トークンで保護された登録をサポートし、乱用対策・審査フローに合わせて運用可（§4）。
- メタデータ標準化: `redirect_uris`、`grant_types`、`token_endpoint_auth_method`、`jwks_uri` などで相互運用性（§3）。
- 登録後管理: `registration_access_token`/`registration_client_uri` により更新・削除等を API で管理（§5, RFC 7592）。

## 留意点（ネイティブ/ローカル・クライアント前提）
- シークレット非保持: 各開発者マシン上の MCP Client は「公開クライアント」として扱い、`client_secret` に依存しない。PKCE の利用を必須化（RFC 7636）。
- リダイレクト URI: ネイティブ向け推奨（RFC 8252）に従い、カスタムスキーム/ループバック（127.0.0.1:random）/プライベート URI を使用。ワイルドカードは避ける。
- 登録エンドポイント保護: 初期アクセス・トークンや `software_statement`（署名 JWT）で乱用防止。加えてレート制御・監査・クォータを実装（RFC 7591 §4, §3, §8）。
- 発行方式・有効期限: `registration_access_token` の保管/ローテーション、失効手段を設計（RFC 7591 §5, RFC 7592）。
- 認証方式ポリシー: `token_endpoint_auth_method` は `none`/`private_key_jwt` 等のサポート方針を明示し、ネイティブには `none+PKCE` や デバイスコード等を優先。
- MCP との責務分離: OAuth は HTTP 接続の保護/認可に用い、MCP 側の `tools/resources/permissions` による最小権限・実行時承認は別途実施。

## 代替検討
- 共有クライアント vs 個別登録: 1つの固定 `client_id` を全ユーザーで共有する方式は配布容易だが、テナント分離・失効単位・利用状況分析で不利。大量配布/個別制御が要件なら DCR が有利。

## 出典（RFC）
- RFC 7591: OAuth 2.0 Dynamic Client Registration Protocol（§1: 動機/背景, §3: メタデータ, §4: 登録エンドポイント/保護, §5: レスポンス/管理連携, §8: セキュリティ）
- RFC 7592: OAuth 2.0 Dynamic Client Registration Management Protocol（登録後の管理 API）
- RFC 6749: The OAuth 2.0 Authorization Framework（基礎概念）
- RFC 7636: Proof Key for Code Exchange by OAuth Public Clients（PKCE）
- RFC 8252: OAuth 2.0 for Native Apps（ネイティブアプリのベストプラクティス）
