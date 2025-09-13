# 実装計画（Go）

## 技術スタック
- 言語: Go 1.22+
- HTTP: 標準 `net/http`（必要なら `github.com/go-chi/chi/v5`）
- JWT/JWK: `github.com/go-jose/go-jose/v4`
- テンプレート: HTML（簡易ログイン/同意）
- ストレージ: インメモリ（map + TTL / time.Timer）

## ディレクトリ構成（提案）
```
cmd/as-poc/main.go          # エントリポイント
internal/http/handlers.go   # 各エンドポイント
internal/oauth/validator.go # PKCE/state/resource 検証
internal/oauth/issuer.go    # トークン発行（JWT）
internal/oauth/store.go     # コード/デバイスコード/セッション in-memory
internal/oauth/clients.go   # 静的クライアント定義
internal/wellknown/meta.go  # RFC 8414/ JWKS 生成
web/templates/*             # login/consent
```

## 静的クライアント（例）
- `client_id`: `mcp-cli-12345`
- 許可 redirect:
  - `http://127.0.0.1/*`（可変ポート許容。実装でホスト=127.0.0.1/::1、パス=/callback 固定などのバリデーション）
  - `myapp://oauth2redirect`
- client_secret: なし（公開クライアント）
- 許可 scope: `openid profile offline_access mcp.read mcp.write`

## エンドポイント実装手順
1) /.well-known/jwks.json
- RSA P-256 等で鍵を生成（開発用は起動時に固定読み込み）。公開 JWK を配信。

2) /.well-known/oauth-authorization-server（RFC 8414）
- issuer/authorize/token/jwks を返却（Device Authorization は任意）。

3) /oauth2/authorize
- 入力検証: response_type, client_id, redirect_uri, state, code_challenge(S256), scope, resource。
- ログイン（ダミー）→ 同意画面 → `code` 発行（TTL 60〜120秒）。`code` と `code_challenge`/`client_id`/`redirect_uri`/`resource` を紐づけ保管。
- 302 で `redirect_uri?code=...&state=...`。

4) /oauth2/token（authorization_code）
- 入力検証: grant_type, code, redirect_uri, client_id, code_verifier, resource。
- `code` を検索し単回使用と期限を確認。`S256(code_verifier)` = `code_challenge` を検証。
- `resource` が省略された場合は `code` に紐づく値を使用。`aud`/`resource` を JWT に反映して署名し返却。

（任意）Device Authorization を実装する場合
- /oauth2/device_authorization と `grant_type=device_code` を追加し、`interval/slow_down` を実装。

## JWT ペイロード例
```
{
  "iss": "https://as.local.test",
  "sub": "user-123",
  "aud": "https://mcp.example.com",
  "exp": 1735680000,
  "iat": 1735676400,
  "scope": "mcp.read mcp.write"
}
```

## エラー設計（例）
- authorize: `invalid_request`, `unauthorized_client`, `access_denied` → エラーページ
- token: 400 `{ error: invalid_grant | invalid_request | unsupported_grant_type }`
- 401/403: `WWW-Authenticate: Bearer error="invalid_token" ...`

## 開発メモ
- `make dev`: `go run ./cmd/as-poc`
- `make fmt`/`make lint`: gofmt, golangci-lint（任意）
- テスト: PKCE 検証、可変ポートの redirect 許容、JWT 検証
