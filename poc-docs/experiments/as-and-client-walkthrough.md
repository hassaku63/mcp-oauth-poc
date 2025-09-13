# 実験手順（OAuth AS PoC × クライアントCLI, 手動コピー方式）

前提
- Go 1.22 以上
- ブラウザが利用可能
- ループバック `127.0.0.1:{任意ポート}` へのアクセスでアドレスバーURLをコピーできる

## 1. ビルド
```
make build
```
生成物:
- `bin/as-poc`（Authorization Server）
- `bin/auth-cli`（クライアント補助CLI）

## 2. Authorization Server を起動
```
make run-as
```
デフォルトの発行者(issuer): `http://localhost:8080`
- メタデータ: `http://localhost:8080/.well-known/oauth-authorization-server`
- JWKS: `http://localhost:8080/.well-known/jwks.json`
- 認可: `http://localhost:8080/oauth2/authorize`
- トークン: `http://localhost:8080/oauth2/token`

別ターミナルに移ります。

## 3. セッション作成（PKCE+state生成）
```
bin/auth-cli new \
  --redirect-uri http://127.0.0.1:53219/callback \
  --resource http://localhost-resource \
  --scope "mcp.read mcp.write" \
  --client-id mcp-cli-12345
```
出力例（抜粋）:
```
{
  "session_id": "sess_...",
  "code_verifier": "...",
  "code_challenge": "...",
  "state": "...",
  "redirect_uri": "http://127.0.0.1:53219/callback",
  "resource": "http://localhost-resource",
  "scope": "mcp.read mcp.write"
}
```
`session_id` を控えます（例: `sess_abc123`）。

## 4. 認可URLを生成
```
bin/auth-cli url \
  --session-id sess_abc123 \
  --authorize-endpoint http://localhost:8080/oauth2/authorize
```
1 行のURLが表示されます。ブラウザで開き、ログイン/同意（PoCは自動許可）後、ループバックにリダイレクトされます。
このときローカル受信は行っていないためブラウザはエラー表示になりますが、アドレスバーにフルURL（`http://127.0.0.1:53219/callback?code=...&state=...`）が表示されます。

## 5. フルURLを貼り付けてトークン交換
```
bin/auth-cli complete \
  --session-id sess_abc123 \
  --token-endpoint http://localhost:8080/oauth2/token \
  --client-id mcp-cli-12345 \
  --callback-url "http://127.0.0.1:53219/callback?code=...&state=..."
```
出力（例）:
```
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "scope": "mcp.read mcp.write"
}
```

## 6. トークンの確認（任意）
- 開発用途でJWTのヘッダ/ペイロードをBase64URLデコードして `aud`/`scope`/`exp` を確認してください（署名検証はAS/リソースサーバ側の責務）。

## トラブルシュート
- 401/403 が出る: `grant_type`/`client_id`/`redirect_uri`/`code_verifier`/`state` を再確認。認可コードのTTL（標準で2分）切れにも注意。
- state mismatch: `session_id` を取り違えていないか確認。`url` と `complete` で同じセッションを使う。
- ブラウザでのURLコピー: クエリパラメータを省略せず、フルURLをそのまま貼り付けること。

## Make の補助ターゲット
```
make help
make auth-new
make auth-url
make auth-complete
```
環境変数 `ISSUER` を上書きすれば、別ポート/別ホストも試せます。
