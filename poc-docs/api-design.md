# API 設計（最低限）

## /.well-known/oauth-authorization-server（RFC 8414）
- GET 200 application/json
```json
{
  "issuer": "https://as.local.test",
  "authorization_endpoint": "https://as.local.test/oauth2/authorize",
  "token_endpoint": "https://as.local.test/oauth2/token",
  "jwks_uri": "https://as.local.test/.well-known/jwks.json",
  "grant_types_supported": ["authorization_code"],
  "code_challenge_methods_supported": ["S256"],
  "token_endpoint_auth_methods_supported": ["none"]
}
```

## /.well-known/jwks.json
- GET 200 application/json（サーバが署名に使う公開鍵）

## /oauth2/authorize（Authorization Endpoint）
- GET（ブラウザ）
- 必須: `response_type=code`, `client_id`, `redirect_uri`, `state`, `code_challenge`, `code_challenge_method=S256`
- 推奨: `scope`, `resource`
- 振る舞い: ログイン→同意→`302 Location: {redirect_uri}?code=...&state=...`。エラー時は `error`, `error_description`。

## /oauth2/token（Token Endpoint）
- POST application/x-www-form-urlencoded
- グラント: `authorization_code`
  - 必須: `grant_type=authorization_code`, `code`, `redirect_uri`, `client_id`, `code_verifier`
  - 推奨: `resource`
- 成功: 200 JSON `{ access_token, token_type="Bearer", expires_in, refresh_token? , scope }`
- 失敗: 400/401 JSON `{ error, error_description }` + `WWW-Authenticate: Bearer error="..."`

（任意）Device Authorization を実装する場合は、`device_authorization_endpoint` と `urn:ietf:params:oauth:grant-type:device_code` を追加する。

## サンプル：WWW-Authenticate（RFC 6750）
```
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Bearer realm="as", error="invalid_token", error_description="expired"
```
