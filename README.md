# README

このリポジトリは、OAuth 2.0 における "Dynamic Client Registration" の仕様理解を目的とした学習用リポジトリです。

以下のことをします。

1. RFC ドキュメントを読み、仕様を理解する
2. (未定) 仕様に基づいて、サンプルコードを実装する (Go 言語)

また、この調査の最終目的は自作の MCP (Model Context Protocol) の認可機構として Dynamic Client Registration を組み込むことです。
本ドキュメント記載時点での MCP 最新仕様 (Version 2025-06-18) では Dynamic Client Registration への対応が明言されています。


## References

- [RFC 7591 - OAuth 2.0 Dynamic Client Registration Protocol](https://datatracker.ietf.org/doc/html/rfc7591)
- [RFC 6749 -  The OAuth 2.0 Authorization Framework](https://datatracker.ietf.org/doc/html/rfc6749)
- [RFC 8414 - OAuth 2.0 Authorization Server Metadata](https://datatracker.ietf.org/doc/html/rfc8414)
- [RFC 9728 - OAuth 2.0 Protected Resource Metadata](https://datatracker.ietf.org/doc/html/rfc9728)
- [Model Context Protocol](https://modelcontextprotocol.io/docs/getting-started/intro)
    - [Overview](https://modelcontextprotocol.io/specification/2025-06-18/basic)
    - [Lifecycle](https://modelcontextprotocol.io/specification/2025-06-18/basic/lifecycle)
    - [Authorization](https://modelcontextprotocol.io/specification/2025-06-18/basic/authorization)
    - [Security Best Practices](https://modelcontextprotocol.io/specification/2025-06-18/basic/security_best_practices)

## 補助資料

- [kura - OAuth/OpenID Connectで実現するMCPのセキュアなアクセス管理](https://speakerdeck.com/kuralab/20250620-openid-technight-mcp-oauth)
- [hi120ki - MCPの認証と認可の現在と未来](https://hi120ki.github.io/ja/blog/posts/20250728/)
- [Hi120ki - MCPの認証と認可 - MCP Meetup Tokyo 2025](https://speakerdeck.com/hi120ki/mcp-authorization)

## What’s Inside
- PoC（Go）: Authorization Server + Protected Resource
  - Auth Code + PKCE（S256）, 手動コピー方式, RFC 6750 準拠の `WWW-Authenticate`
  - 実装: `cmd/as-poc`, `internal/server/*`
- PoC: クライアント補助 CLI（`new`/`url`/`complete`）
  - PKCE・state生成、認可URL生成、貼付URL解析＋トークン交換
  - 実装: `cmd/auth-cli`, `internal/client/*`
- 設計ドキュメント
  - MCP 認可の要点とフロー: `docs/mcp-authorization-notes.md`
  - 役割別プロトコル要件（MUST/SHOULD + RFC）: `docs/mcp-role-protocol-requirements.md`
  - 実装計画: `poc-docs/oauth-as/*`, `poc-docs/oauth-client/*`, `poc-docs/protected-resource/*`

## Quick Start
- Build: `make build`（生成: `bin/as-poc`, `bin/auth-cli`）
- Run AS: `make run-as`（issuer: `http://localhost:8080`）
- Walkthrough: `poc-docs/experiments/as-and-client-walkthrough.md`
  - `auth-cli new` → `auth-cli url` → ブラウザ同意 → callback URL をコピー → `auth-cli complete`
  - 取得トークンで `GET /resource/echo` を `Authorization: Bearer` で呼ぶ

## MCP × OAuth Roles（要点）
- MCP Client（OAuth Client）
  - MUST: Auth Code + PKCE（S256）＋ランダム state（RFC 6749/7636）; ネイティブは RFC 8252（loopback/custom scheme）
- Authorization Server（AS）
  - MUST: RFC 8414 メタデータ, PKCE対応の authorize/token, RFC 6750 の挑戦/エラー, JWKS 提供
- MCP Server（Resource Server）
  - MUST: RFC 9728（Protected Resource Metadata）公開, `iss/aud/exp`/署名検証, 401/403 は RFC 6750 準拠（`WWW-Authenticate`）
  - SHOULD: スコープ→操作の対応表で最小権限（例: `mcp.read` 必須）
