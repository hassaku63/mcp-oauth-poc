# 要件（MUST/SHOULD）— Protected Resource 検証

MUST（最低限）
- `GET /resource/echo` を提供し、Bearer トークン必須とする。
- JWT 検証: 署名（RS256, ASの鍵）、`iss` 一致、`exp`/`iat` 妥当、`aud`（= resource識別子）一致、`scope` に `mcp.read` を含むこと。
- エラー応答（RFC 6750）:
  - 401 Unauthorized + `WWW-Authenticate: Bearer error="invalid_token"`（欠落/不正/期限切れ/issuer不一致/audience不一致/署名不正）
  - 403 Forbidden + `WWW-Authenticate: Bearer error="insufficient_scope", scope="mcp.read"`

SHOULD（推奨）
- `WWW-Authenticate` の `realm`（例: `resource`）や `error_description` をわかりやすく付与。
- 設定で期待 `aud`（resource識別子）を一元管理（例: `http://localhost-resource`）。
- ログにはトークン本体を出さず、検証結果・理由のみ記録。

OPTIONAL（任意）
- `POST /resource/echo` を追加し、`mcp.write` が必要なパスで 403 検証をしやすくする。
- JWKS フェッチ/キャッシュ（将来 AS/Resource を分離する前提整備）。
