# PoC: MCP クライアント補助ツール（Authorization Code + PKCE, 手動コピー方式）

前提: MCP Server の Transport に HTTP/WS を用いる構成で、Authorization Code + PKCE を前提とします。

目的
- MCP Client（CLI/デスクトップ）が Authorization Code + PKCE を用いて AS と連携するための最小ツール群を提供する。
- ループバック受信なし（手動コピー方式）を前提とし、UXを損なわずに安全性を確保する。

到達目標（MVP）
- PKCE 生成（`code_verifier` / `code_challenge(S256)`）と `state` 生成・保存（TTL/単回利用）。
- 認可URL生成（`response_type=code`, `client_id`, `redirect_uri`, `scope`, `state`, `code_challenge`, `resource`）。
- 手動貼付 URL の解析（`code`/`state` 抽出・照合）とトークン交換。
- メタデータ自動解決（任意）: RFC 9728（MCP Server → AS発見）→ RFC 8414（AS メタデータ）。

非目標（この PoC では扱わない）
- ループバック自動受信サーバ、カスタムスキーム・ハンドラ。
- Device Authorization（必要になれば任意で追加）。

関連ドキュメント
- requirements.md: 必須/推奨要件
- api-design.md: CLI コマンド仕様
- sequence-diagrams.md: フロー図
- impl-plan.md: 実装手順（Go）
- security-notes.md: セキュリティ観点
