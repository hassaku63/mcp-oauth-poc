# MCP Server における RFC 9728 要件 — 報告

## 結論
- 必須: HTTP Transport を用いる MCP Server は、自身に紐づく Authorization Server の場所を示すために、OAuth 2.0 Protected Resource Metadata（RFC 9728）を公開しなければならない（MUST）。
- 前提: この要件は Dynamic Client Registration（DCR）を使わない場合でも変わらない（独立の要件）。

## 理由
- DCR は OAuth クライアントのプロビジョニング方式に関する仕様であり、どの認可サーバを使うかの“発見（discovery）”とは別問題。
- MCP Client は、どの Authorization Server でトークンを取得すべきかを標準的に知る必要があるため、MCP Server 側の RFC 9728 公開が必須となる。

## 公開すべき内容（最小）
- 配置: `/.well-known/oauth-protected-resource`
- 必須フィールド:
  - `resource`: MCP Server を一意に表す保護リソース識別子（URL/URN）。
  - `authorization_servers`: 利用可能な Authorization Server（発行者）の URL 配列。

```json
{
  "resource": "https://mcp.example.com",
  "authorization_servers": [
    "https://auth.example.com"
  ]
}
```

## クライアント側の後続手順（推奨）
- RFC 8414: 列挙された各 Authorization Server のメタデータを取得して、トークン/認可エンドポイント等を解決する。
- RFC 8707: トークン取得時に `resource` を指定し、MCP Server 向けにスコープされたアクセストークンを要求する。
- RFC 6750: Bearer トークン提示と 401/403 時の `WWW-Authenticate` 取り扱いに準拠する。

## 備考
- 複数の Authorization Server の提示も可能で、クライアントはポリシーに応じて選択できる。
- `stdio` などローカル接続では OAuth による接続認可を用いない構成もあるが、HTTP/WS を用いる場合は本要件が適用される。
