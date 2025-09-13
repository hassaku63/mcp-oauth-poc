# 要件（MUST/SHOULD）— MCP クライアント補助ツール

MUST（最低限）
- PKCE 生成: `code_verifier`（43–128 文字, base64url, 無パディング）, `code_challenge=S256(verifier)` を生成。
- state 生成/検証: 毎回ランダム生成し、保存→手動貼付URLからの `state` と一致を検証（TTL/単回利用）。
- 認可URL生成: `response_type=code`, `client_id`, `redirect_uri`, `scope`, `state`, `code_challenge`, `code_challenge_method=S256`, `resource` をURLエンコードして構築。
- 手動貼付URL解析: `code`/`state` を抽出し、保存済みの `state` と一致確認。
- トークン交換: `grant_type=authorization_code`, `code`, `redirect_uri`, `client_id`, `code_verifier`, `resource` を `application/x-www-form-urlencoded` で送信。

SHOULD（推奨）
- メタデータ解決: RFC 9728（`/.well-known/oauth-protected-resource`）→ RFC 8414（AS metadata）で `authorization_endpoint`/`token_endpoint` を自動取得。
- エラー整備: `invalid_grant`/`invalid_request`/ネットワークエラーのわかりやすい表示。
- クリップボード連携: 認可URLのコピー/貼付支援、ブラウザ起動（環境依存）。
- ログ衛生: `code_verifier`/トークンをログに出さない。最低限の診断情報のみ記録。

OPTIONAL（任意）
- ループバック自動受信（RFC 8252）への拡張、Device Authorization（RFC 8628）。
