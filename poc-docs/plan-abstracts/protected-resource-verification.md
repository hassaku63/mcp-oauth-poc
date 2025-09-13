# 次フェーズ計画: Protected Resource へのアクセス検証（OAuth 2.0 Bearer）

## 目的
- 取得済みアクセストークンで保護リソースにアクセスできることを、RFC 6750（Bearer Token Usage）に沿って検証する。
- 署名・期限・発行者・オーディエンス（resource/aud）・スコープをサーバ側で厳格に検証し、適切な 401/403 と `WWW-Authenticate` を返す。

## スコープ（MVP）
- エンドポイント追加: `GET /resource/echo`
  - 認可要件: スコープ `mcp.read` を要求
  - 動作: `Authorization: Bearer <access_token>` を検証し、ペイロード（sub, scope, aud など）をJSONで返却
- トークン検証
  - 署名: RS256（現行ASの鍵に一致）
  - `iss`: `cfg.Issuer` と一致
  - `aud`（または `resource` 相当）: 期待値（例: `http://localhost-resource`）と一致
  - 有効期限: `exp` > now、`iat` は許容スキュー内
  - スコープ: 要求された操作に必要な権限を含む（`mcp.read`）
- エラーハンドリング（RFC 6750）
  - 401 + `WWW-Authenticate: Bearer realm="resource", error="invalid_token", error_description="..."`
  - 403 + `WWW-Authenticate: Bearer realm="resource", error="insufficient_scope", scope="mcp.read"`

## 実装タスク
- server 層
  - [ ] `internal/server/resource/echo.go`: `GET /resource/echo` 実装
  - [ ] `internal/server/http/middleware.go`: Bearer 抽出/エラー応答ヘルパ（WWW-Authenticate組み立て）
  - [ ] `cmd/as-poc/main.go`: ルーティング追加（`/resource/echo`）
- 検証ロジック
  - [ ] `internal/server/oauth/validate_token.go`（新規）: JWT 解析・署名検証・クレーム検証（iss/aud/exp/iat/scope）
  - [ ] 設定: 期待 `aud`（= Resource Identifier）を `Config` に追加（例: `ResourceID string`）。現時点のPoC既定値は `http://localhost-resource`。
- ドキュメント/手順
  - [ ] `poc-docs/experiments/as-and-client-walkthrough.md`: 取得トークンで `/resource/echo` を curl する手順を追加
  - 例: `curl -sS http://localhost:8080/resource/echo -H "Authorization: Bearer $ACCESS_TOKEN"`

## 受け入れ基準
- 正しいトークンで 200、本文に JWT クレームの要点（sub/scope/aud/exp）が含まれる
- トークン欠落/形式不正/署名不正/期限切れ/issuer不一致/audience不一致 → 401 + `WWW-Authenticate: Bearer error="invalid_token"`
- スコープ不足（例: `mcp.read` が無い）→ 403 + `WWW-Authenticate: Bearer error="insufficient_scope", scope="mcp.read"`

## 追加提案（任意強化）
- JWKS 連携: リソースサーバが AS の `jwks_uri` から鍵を取得・キャッシュ（将来 AS/Resource を別プロセス化する前提整備）
- Scope ベースの操作拡張: `POST /resource/echo`（`mcp.write` 必須）で 403 の検証パスを容易に再現
- CLI 拡張: `auth-cli call --url http://localhost:8080/resource/echo --access-token ...`（Authorization ヘッダ付与の簡易ヘルパ）
- リフレッシュトークン（任意）: `refresh_token` 付与と `grant_type=refresh_token` の追加で長期セッション試験
- 監査/メトリクス: 失敗理由の分類ログ、レイテンシ/成功率の計測
- RFC 9728（参考）: 本PoCではAS側要件だが、将来はリソース側でも `.well-known/oauth-protected-resource` を返し、ASロケーションを自己記述

## サンプル応答（参考）
- 401（期限切れ）
```
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Bearer realm="resource", error="invalid_token", error_description="expired"
```
- 403（スコープ不足）
```
HTTP/1.1 403 Forbidden
WWW-Authenticate: Bearer realm="resource", error="insufficient_scope", scope="mcp.read"
```
