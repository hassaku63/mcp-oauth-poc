# PoC: Protected Resource 検証（OAuth 2.0 Bearer + JWT）

目的
- 取得済みアクセストークンで保護リソースにアクセスできることを検証し、RFC 6750（Bearer）準拠の挑戦/エラー応答と、JWT クレーム（iss/aud/exp/scope 等）の検証を行う。

到達目標（MVP）
- `GET /resource/echo` を追加し、`Authorization: Bearer <token>` を検証して 200/401/403 を適切に返す。
- スコープ `mcp.read` 必須。トークンの `iss`/署名/有効期限/`aud`（= resource）/`scope` を検証。
- `WWW-Authenticate` 応答を RFC 6750 に沿って整備（`invalid_token`/`insufficient_scope`）。

関連ドキュメント
- requirements.md: MUST/SHOULD 要件
- api-design.md: エンドポイント仕様
- sequence-diagrams.md: フロー図
- impl-plan.md: Go 実装手順
- security-notes.md: セキュリティ観点
