# 実装計画（Go）— Protected Resource

## 変更点（server 側）
- ルーティング: `cmd/as-poc/main.go` に `GET /resource/echo` を追加。
- 新規: `internal/server/resource/echo.go` — ハンドラ実装。
- 新規: `internal/server/http/middleware.go` — Bearer 抽出/エラー応答/WWW-Authenticateヘルパ。
- 新規: `internal/server/oauth/validate_token.go` — JWT 解析・署名検証・iss/aud/exp/iat/scope 検証。
- 既存: `internal/server/oauth/types.go` の `Config` に `ResourceID string` を追加（例: `http://localhost-resource`）。

## 検証仕様
- 署名: RS256（現行 `Config.PrivateKey` で検証可）。
- iss: `Config.Issuer` と一致。
- exp/iat: 現在時刻と比較（スキュー許容）。
- aud: `Config.ResourceID` と一致（無ければ 401 `invalid_token`）。
- scope: `mcp.read` を含まない場合 403 `insufficient_scope`。

## ハンドラの振る舞い
- 成功: 200 JSON（sub/aud/scope/exp の要点を返す）
- 失敗: RFC 6750 に準拠した `WWW-Authenticate`
  - 401: `Bearer realm="resource", error="invalid_token", error_description="..."`
  - 403: `Bearer realm="resource", error="insufficient_scope", scope="mcp.read"`

## 手順
1) `types.go` に `ResourceID` を追加し、`cmd/as-poc/main.go` でデフォルトを設定。
2) `validate_token.go` を作成（JWTパース/署名検証/クレーム検証）。
3) `middleware.go` を作成（Authorizationヘッダ抽出、エラー/WWW-Authenticateを一元化）。
4) `resource/echo.go` を実装し、ミドルウェア＋バリデータを呼び出し。
5) `main.go` にルーティングを追加。
6) 実験手順ドキュメント（experiments）に `/resource/echo` への curl 例を追記。

## 受け入れ基準
- 正しいトークン→200、必要クレームが本文に含まれる。
- 欠落/不正トークン→401 `invalid_token`。
- スコープ不足→403 `insufficient_scope`。
