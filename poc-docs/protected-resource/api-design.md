# API 設計（Protected Resource）

## GET /resource/echo
- 要件: `Authorization: Bearer <access_token>`（`mcp.read` 必須）
- 正常系: 200 application/json
```json
{
  "sub": "user-123",
  "aud": "http://localhost-resource",
  "scope": "mcp.read mcp.write",
  "exp": 1735680000
}
```
- エラー系（RFC 6750）
  - 欠落/不正/期限切れなど: 401 +
```
WWW-Authenticate: Bearer realm="resource", error="invalid_token", error_description="..."
```
  - スコープ不足: 403 +
```
WWW-Authenticate: Bearer realm="resource", error="insufficient_scope", scope="mcp.read"
```

## （任意）POST /resource/echo
- 要件: `mcp.write` を要求
- 動作: POST された JSON をそのまま返す（エコー）。
