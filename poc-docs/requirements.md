# 要件（MUST/SHOULD）

MUST（最低限実装）
- Authorization Endpoint: `GET /oauth2/authorize`（Authorization Code + PKCE S256、`state` 必須）。
- Token Endpoint: `POST /oauth2/token`（authorization_code）。
- Authorization Server Metadata（RFC 8414）: `GET /.well-known/oauth-authorization-server`。
- JWKS（RFC 7517）: `GET /.well-known/jwks.json`。
- Bearer Token Usage（RFC 6750）: エラー応答に `WWW-Authenticate` を含む。
- Resource Indicators（RFC 8707）: `resource` パラメータを受理・検証し、発行トークンの宛先に反映（`aud`/カスタム claim）。
- PKCE Validation: `code_verifier` を S256 で検証。
- 静的クライアント: `client_id` と許可 `redirect_uri` の静的登録（公開クライアントとして `client_secret` 不要）。

SHOULD（推奨）
- ループバックリダイレクト: 127.0.0.1/::1 の可変ポートを許容（RFC 8252）。
- ログイン/同意 UI: PoC 用の簡易ログインとスコープ同意 UI。
- トークン形式: JWT アクセストークン（RS256）で発行。`aud`/`exp`/`iat`/`scope`/`cnf?` 等。
- メモリ永続: 認可コードの TTL 管理（in-memory）。
- レート制御/CSRF/XSS 対策（最低限）。

OPTIONAL（任意/将来）
- Device Authorization Endpoint（RFC 8628）と `urn:ietf:params:oauth:grant-type:device_code` 対応。

OUT（非対応）
- Dynamic Client Registration（RFC 7591/7592）。
- Introspection/Revocation（PoC では省略可）。
- mTLS / JAR / Pushed Authorization Requests 等。
