# CLI コマンド設計（最小）

コマンド群（例: `auth-cli`）

## new
- 説明: PKCE と state を生成し、セッションIDとともに保存。JSON を返す。
- 入力: `--redirect-uri`, `--resource`, `--scope`, `--client-id`
- 出力(JSON例):
```json
{
  "session_id": "sess_abc123",
  "code_verifier": "...",
  "code_challenge": "...",
  "state": "...",
  "redirect_uri": "http://127.0.0.1:53219/callback",
  "resource": "https://mcp.example.com",
  "scope": "mcp.read mcp.write"
}
```

## url
- 説明: `new` の結果と AS メタデータから認可URLを生成して出力（オプションで既定ブラウザを起動）。
- 入力: `--session-id`, `--authorize-endpoint`（又は `--as`/`--issuer` から自動解決）
- 出力: 認可URL（1行）

## complete
- 説明: ユーザーが貼り付けたフル callback URL から `code`/`state` を抽出し、state検証→トークン交換。
- 入力: `--session-id`, `--token-endpoint`（又は `--as`/`--issuer` から自動解決）, `--client-id`, `--callback-url`
- 出力(JSON例): `{ "access_token": "...", "token_type": "Bearer", "expires_in": 3600, "scope": "..." }`

## discover（任意）
- 説明: MCP Server の RFC 9728 を取得して AS を列挙、RFC 8414 でエンドポイントを解決。
- 入力: `--resource-metadata-url`（例: `https://mcp.example.com/.well-known/oauth-protected-resource`）
- 出力: `authorization_endpoint`, `token_endpoint`, `issuer`
