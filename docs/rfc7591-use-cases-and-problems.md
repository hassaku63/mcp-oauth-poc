# RFC 7591 Dynamic Client Registration — 課題とユースケース

## 解決する課題（What it solves）
- 標準的な登録手段の欠如: OAuth 2.0（RFC 6749）ではクライアント登録は前提だが手段は規定されず、多くが手作業・事前合意（out-of-band）で非効率かつ非互換だった問題を解消し、動的登録の標準プロトコルを定義。
  - 出典: RFC 7591 §1（Introduction）
- スケールと自動化: 多数のクライアント/テナント/配布形態に対し、管理者の手作業なしで `client_id` 等を払い出し可能にして、SaaS/モバイル/配布エコシステムでの展開を自動化。
  - 出典: RFC 7591 §1、§4（Client Registration Endpoint）
- ポリシーに応じた登録制御: オープン登録（誰でも登録）から、初期アクセス・トークンによる保護付き登録まで幅広いポリシーを許容。
  - 出典: RFC 7591 §4（Registration Endpoint の保護方針に関する言及）
- メタデータの相互運用性: クライアントの性質・挙動を示す標準メタデータ（`redirect_uris`、`grant_types`、`token_endpoint_auth_method`、`jwks_uri` など）を規定して相互運用性を向上。
  - 出典: RFC 7591 §3（Client Metadata）
- ソフトウェア・ステートメントによる信頼移送: 配布元/レジストリが署名した `software_statement`（JWT）でメタデータ/身元を伝達でき、登録時の審査・信頼モデルに組み込める。
  - 出典: RFC 7591 §3（`software_statement` 定義）

## 想定ユースケース（Use cases）
- モバイル/ネイティブ/デスクトップ・アプリの配布: アプリが実行時に登録して `client_id`（必要なら `client_secret` 以外の認証方式）を取得。
  - 出典: RFC 7591 §1、§3（公開クライアント/認証方法のバリエーション）
- SaaS のテナント毎クライアント: 管理 UI や自動プロビジョニングから動的にクライアントを作成し、テナント分離や構成管理を容易化。
  - 出典: RFC 7591 §1（自動化・大規模展開の動機）、§4
- マーケットプレイス/レジストリ経由の配布: ベンダが署名した `software_statement` を用い、AS 側がポリシーに沿って安全に受け入れ。
  - 出典: RFC 7591 §3（`software_statement`）
- オープン・エコシステムでのセルフサービス登録: 事前の人手審査なしで、ポリシー/レート制御/審査後の無効化等を前提にセルフサービス登録を許す。
  - 出典: RFC 7591 §1（オープン登録の言及）、§4（エンドポイント保護）
- 後続のクライアント管理（読み取り/更新/削除）: 登録後の構成 API は別仕様（RFC 7592）で定義され、`registration_access_token`/`registration_client_uri` を通じて管理操作を実施。
  - 出典: RFC 7591 §5（レスポンス項目の定義と 7592 への参照）、RFC 7592（Management）

## メモ（仕様上のキーポイント）
- レジストレーション・エンドポイント: HTTP `POST` でメタデータを送信し、成功時に `client_id` 等と共に（必要に応じて）`registration_access_token` と `registration_client_uri` を受領。
  - 出典: RFC 7591 §4, §5
- 保護モデル: AS のポリシーにより、無認証（オープン）または「初期アクセス・トークン」を要求する保護付き登録のいずれも可。
  - 出典: RFC 7591 §4
- セキュリティ考慮: メタデータの検証、HTTPS 要求、発行する認証方式の制限、レート制御/監査等。
  - 出典: RFC 7591 §8（Security Considerations）

## 参考 RFC
- RFC 7591: OAuth 2.0 Dynamic Client Registration Protocol（特に §1, §3, §4, §5, §8）
- RFC 7592: OAuth 2.0 Dynamic Client Registration Management Protocol（登録後の管理）
- RFC 6749: The OAuth 2.0 Authorization Framework（クライアント登録の前提と概念）
