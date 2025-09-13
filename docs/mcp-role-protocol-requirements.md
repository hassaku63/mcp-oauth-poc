# MCP ロール別プロトコル要件（要約）

参照元（例）
- MCP Basic Spec — Authorization (2025-06-18): https://modelcontextprotocol.io/specification/2025-06-18/basic/authorization
- スライド: https://speakerdeck.com/kuralab/20250620-openid-technight-mcp-oauth?slide=17
- 関連 RFC: 6749, 6750, 7636, 8252, 8414, 8707, 9728, 7591/7592, 8628（任意）

## MCP Client（OAuth Client）
- MUST [RFC 6749][RFC 7636]: Authorization Code + PKCE（S256）を使用
- MUST [RFC 6749]: `state` を毎回ランダム生成・検証（CSRF/コード差し替え防止）
- MUST [RFC 8252]: ネイティブ/CLI は RFC 8252 に従う（ループバック 127.0.0.1/::1 可変ポート、またはカスタムスキーム）
- MUST [RFC 6750]: TLS を使用し、`Authorization: Bearer` で提示
- SHOULD [RFC 9728][RFC 8414]: Protected Resource Metadata → AS Metadata でエンドポイント発見
- SHOULD [RFC 8707]: トークン取得時に `resource` を指定（MCP Server 宛にスコープ）
- SHOULD: トークン/`code_verifier`/`state` の安全管理（短命・単回利用・非ログ）
- MAY [RFC 8628]: Device Authorization をフォールバックとして実装
- MUST NOT [RFC 6749][RFC 8252]: Implicit フロー/固定 `state`/ネイティブでの `client_secret` 依存

## Authorization Server（AS）
- MUST [RFC 8414]: Authorization Server Metadata を公開（`/.well-known/oauth-authorization-server`）
- MUST [RFC 6749][RFC 7636]: Authorization Code + PKCE（S256）と token endpoint を実装
- MUST [RFC 6750]: Bearer のエラー/挑戦を準拠表示
- MUST [RFC 7517][RFC 8414]: JWKS（`jwks_uri`）で公開鍵を提供
- SHOULD [RFC 8252]: ループバック可変ポート/カスタムスキーム（ネイティブ対応）
- SHOULD [RFC 8707]: Resource Indicators をサポート（`aud`/`resource` 反映）
- MAY [RFC 7591][RFC 7592]: Dynamic Client Registration / Management を提供
- MAY [RFC 8628]: Device Authorization を提供

## MCP Server（Resource Server）
- MUST [RFC 9728]: Protected Resource Metadata を公開（`/.well-known/oauth-protected-resource`）
- MUST [RFC 7517][RFC 8414][RFC 8707]: アクセストークン検証（署名/JWKS, `iss`, `aud`/`resource`, `exp`/`nbf`/`iat`）
- MUST [RFC 6750]: 失敗時は `WWW-Authenticate: Bearer ...`（401/403）
- SHOULD [RFC 6749]: スコープ→操作のマッピングを定義し各ルートに適用
- SHOULD: TLS を強制（必要に応じて mTLS）
- SHOULD [RFC 8414]: 複数AS対応（issuer許可リスト/JWKSキャッシュ/kidローテ）
- MUST（MCP運用）: 公開surface最小化・危険操作は Permissions でユーザー承認（MCP Spec）
- NOTE: stdio では接続レイヤの OAuth を使わない構成もあるが、HTTP/WS では本要件が適用

## 位置付け（まとめ）
- クライアントは「AS から適切にトークンを取得」し、サーバは「RS として標準に従い検証・応答」する。MCP はその上で、公開機能の最小化と実行時 Permissions により、操作粒度の安全性を高める二層設計となる。

---
参考RFC（略号→正式名）
- [RFC 6749]: The OAuth 2.0 Authorization Framework
- [RFC 6750]: The OAuth 2.0 Authorization Framework: Bearer Token Usage
- [RFC 7636]: Proof Key for Code Exchange (PKCE) by OAuth Public Clients
- [RFC 8252]: OAuth 2.0 for Native Apps
- [RFC 8414]: OAuth 2.0 Authorization Server Metadata
- [RFC 8707]: Resource Indicators for OAuth 2.0
- [RFC 9728]: OAuth 2.0 Protected Resource Metadata
- [RFC 7517]: JSON Web Key (JWK)
- [RFC 7591]: OAuth 2.0 Dynamic Client Registration Protocol
- [RFC 7592]: OAuth 2.0 Dynamic Client Registration Management Protocol
- [RFC 8628]: OAuth 2.0 Device Authorization Grant
