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
