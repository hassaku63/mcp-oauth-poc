# PoC: 最低限の OAuth 2.0 Authorization Server（DCR なし）

前提: MCP Server の Transport に HTTP/WS を用いる前提での検証です。

目的
- MCP Server の HTTP 認可前提に合わせ、最小限の Authorization Server 機能を持つテスト用 AS を Go で実装する。
- Dynamic Client Registration（RFC 7591/7592）は除外。クライアントは静的に定義。

到達目標（MVP）
- Authorization Code + PKCE（S256）対応（公開クライアント前提）。
- 手動コピー方式のサポート（ループバック URI にリダイレクトし、ユーザーがブラウザのフル URL を CLI に貼り付ける運用）。
- OAuth 2.0 Bearer Token（RFC 6750）でアクセストークン発行・提示。
- Authorization Server Metadata（RFC 8414）と JWKS（JWK Set）提供。
- Protected Resource Metadata（RFC 9728）のサンプル（AS から見た関連 MCP リソースの提示は任意。主に MCP Server 側実装だが相互試験用に用意）。
- Resource Indicators（RFC 8707）の `resource` パラメータを受理し、発行トークンに対象を反映（`aud`/`resource`）。

非目標（この PoC では扱わない）
- 動的クライアント登録（RFC 7591/7592）。
- OIDC ID Token/ユーザープロファイルの本格実装（必要ならダミーのユーザーログイン画面）。
- mTLS、JAR/DPoP、JWT ログインなど拡張仕様。

任意機能（必要に応じて）
- Device Authorization Grant（RFC 8628）

関連ドキュメント
- requirements.md: 必須/推奨要件
- api-design.md: エンドポイント仕様
- sequence-diagrams.md: フローのシーケンス図
- impl-plan.md: Go 実装手順
- security-notes.md: セキュリティ観点
- 次フェーズ計画: ../plans/oauth-as-next-steps.md
