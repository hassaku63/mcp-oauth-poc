# 実装計画（Go, 単体バイナリ）

## スタック
- 言語: Go 1.22+
- 乱数/ハッシュ: `crypto/rand`, `crypto/sha256`
- HTTP: `net/http`
- 保存: 一時ディレクトリに JSON（セッション毎） or メモリ（プロセス内）

## ディレクトリ案
```
cmd/auth-cli/main.go        # サブコマンド: new/url/complete/discover
internal/pkce/pkce.go       # PKCE生成
internal/state/state.go     # state生成/保存/検証
internal/authurl/build.go   # 認可URL生成
internal/callback/parse.go  # 貼付URL解析
internal/token/exchange.go  # トークン交換
internal/discover/discover.go # RFC 9728/8414 解決（任意）
```

## 手順
1) `new`: PKCE + state 生成、セッションIDを払い出し、JSON保存（TTL/単回利用フラグ）。
2) `url`: セッション読み込み→AS メタデータ解決（任意）→認可URL生成・表示（`--open`でブラウザ起動）。
3) `complete`: `--callback-url` を解析→`state`照合→トークン交換（HTTP 400/401のハンドリング）。
4) `discover`（任意）: RFC 9728 → RFC 8414 でエンドポイントを導出し、JSONで保存/出力。

## テスト
- PKCE S256 正当性、state 単回利用、URL生成のエンコード、貼付URL解析の堅牢性。
- トークン交換の必須パラメータ、エラー応答の扱い（`invalid_grant`等）。
