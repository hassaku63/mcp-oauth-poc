# MCP ロール別プロトコル要件（要約）

参照元（例）
- MCP Basic Spec — Authorization (2025-06-18): https://modelcontextprotocol.io/specification/2025-06-18/basic/authorization
- スライド: https://speakerdeck.com/kuralab/20250620-openid-technight-mcp-oauth?slide=17
- 関連 RFC:
    - RFC 6749: The OAuth 2.0 Authorization Framework
    - RFC 6750: The OAuth 2.0 Authorization Framework: Bearer Token Usage
    - RFC 7636: Proof Key for Code Exchange (PKCE) by OAuth Public Clients
    - RFC 8252: OAuth 2.0 for Native Apps
    - RFC 8414: OAuth 2.0 Authorization Server Metadata
    - RFC 8707: Resource Indicators for OAuth 2.0
    - RFC 9728: OAuth 2.0 Protected Resource Metadata
    - RFC 7591: OAuth 2.0 Dynamic Client Registration Protocol（任意）
    - RFC 7592: OAuth 2.0 Dynamic Client Registration Management Protocol（任意）
- RFC 8628: OAuth 2.0 Device Authorization Grant（任意）

前提: 本要件は MCP Server の Transport に HTTP/WS を想定して整理しています（stdio 等では接続レイヤの OAuth 要件は必須ではありません）。

## MCP Client（OAuth Client）

| Level | Requirement | RFC |
| - | - | - |
| MUST | Authorization Code + PKCE（S256）を使用 | RFC 6749, RFC 7636 |
| MUST | `state` を毎回ランダム生成・検証（CSRF/コード差し替え防止） | RFC 6749 |
| MUST | ネイティブ/CLI は RFC 8252（loopback 可変ポート or custom scheme）に従う | RFC 8252 |
| MUST | TLS 使用・`Authorization: Bearer` 提示 | RFC 6750 |
| SHOULD | Protected Resource Metadata → AS Metadata でエンドポイント発見 | RFC 9728, RFC 8414 |
| SHOULD | トークン取得時に `resource` 指定（MCP Server 宛にスコープ） | RFC 8707 |
| SHOULD | トークン/`code_verifier`/`state` の安全管理（短命・単回・非ログ） |  |
| MAY | Device Authorization をフォールバックとして実装 | RFC 8628 |
| MUST NOT | Implicit/固定 `state`/ネイティブでの `client_secret` 依存 | RFC 6749, RFC 8252 |

## Authorization Server（AS）

| Level | Requirement | RFC |
| - | - | - |
| MUST | Authorization Server Metadata を公開（`/.well-known/oauth-authorization-server`） | RFC 8414 |
| MUST | Authorization Code + PKCE（S256）と token endpoint を実装 | RFC 6749, RFC 7636 |
| MUST | Bearer のエラー/挑戦を準拠表示 | RFC 6750 |
| MUST | JWKS（`jwks_uri`）で公開鍵を提供 | RFC 7517, RFC 8414 |
| SHOULD | ループバック可変ポート/カスタムスキーム（ネイティブ対応） | RFC 8252 |
| SHOULD | Resource Indicators をサポート（`aud`/`resource` 反映） | RFC 8707 |
| MAY | Dynamic Client Registration / Management を提供 | RFC 7591, RFC 7592 |
| MAY | Device Authorization を提供 | RFC 8628 |

## MCP Server（Resource Server）

| Level | Requirement | RFC |
| - | - | - |
| MUST | Protected Resource Metadata を公開（`/.well-known/oauth-protected-resource`） | RFC 9728 |
| MUST | アクセストークン検証（署名/JWKS, `iss`, `aud`/`resource`, `exp`/`nbf`/`iat`） | RFC 7517, RFC 8414, RFC 8707 |
| MUST | 失敗時は `WWW-Authenticate: Bearer ...`（401/403） | RFC 6750 |
| SHOULD | スコープ→操作のマッピングを定義し各ルートに適用 | RFC 6749 |
| SHOULD | TLS を強制（必要に応じて mTLS） |  |
| SHOULD | 複数AS対応（issuer許可リスト/JWKSキャッシュ/kidローテ） | RFC 8414 |
| MUST（MCP運用） | 公開surface最小化・危険操作は Permissions でユーザー承認（MCP Spec） |  |

NOTE: stdio では接続レイヤの OAuth を使わない構成もあるが、HTTP/WS では本要件が適用

## 位置付け（まとめ）
- クライアントは「AS から適切にトークンを取得」し、サーバは「RS として標準に従い検証・応答」する。MCP はその上で、公開機能の最小化と実行時 Permissions により、操作粒度の安全性を高める二層設計となる。
