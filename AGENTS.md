# Repository Guidelines

このリポジトリは「OAuth 2.0 Dynamic Client Registration (RFC 7591)」の仕様理解を目的とした学習用です。

## Research

`README.md` に記載の RFC ドキュメントを読み、仕様を理解する。

ユーザーの調査リクエストに対する回答は、`docs/*.md` に出力すること。

繰り返し同一 URL の情報への参照が発生する場合は、`docs/cache/*` 以下にキャッシュしても良いものとする。

## 実装ガイド

現時点では、実装は未計画。ドキュメントの読解を優先。

PoC 実装を行う場合は、`poc-docs/<FEATURE_TITLE>/*.md` のパス規則に従うように、設計・実装計画・セキュリティノート等を記載すること。
実際に詳細計画を行うかどうか未定の、概要レベルの記述は `poc-docs/plan-abstract/*.md` にまとめる。
