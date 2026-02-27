# 現場で使えるGo言語実践テクニック — サンプルコード

**設計パターン・並行処理・パフォーマンス最適化**

書籍「現場で使えるGo言語実践テクニック」で紹介するコードサンプルの完全版リポジトリです。

## 動作環境

- **Go 1.22 以上**（一部の章は Go 1.24+ を推奨）
- 付録A のサンプルは Go 1.25/1.26 が必要

## リポジトリ構成

各章のサンプルコードは独立した `go.mod` を持ち、個別にビルド・テスト可能です。

```
go-textbook-advanced/
├── ch01-idiomatic-go/         第1章: Goらしいコードと型の設計
├── ch02-generics/             第2章: ジェネリクスとイテレータ
├── ch03-error-handling/       第3章: エラー処理の設計
├── ch04-design-patterns/      第4章: 設計パターン
├── ch05-concurrency/          第5章: 並行処理パターン
├── ch06-package-design/       第6章: パッケージ設計とモジュール運用
├── ch07-testing/              第7章: テスト戦略
├── ch08-performance/          第8章: パフォーマンス最適化
├── ch09-observability/        第9章: 構造化ログとオブザーバビリティ
├── ch10-database-api/         第10章: データベースとAPI設計
├── ch11-production/           第11章: プロダクション運用
├── ch12-grpc-microservice/    第12章: 実践 gRPCマイクロサービス
├── appendix-a-go126/          付録A: Go 1.25/1.26 新機能
└── appendix-b-review/         付録B: コードレビューチェックリスト
```

## 使い方

### 個別の章を実行

```bash
cd ch01-idiomatic-go
go run .
```

### テストの実行

```bash
cd ch07-testing
go test ./...
```

### ベンチマークの実行

```bash
cd ch08-performance
go test -bench=. -benchmem
```

## 章別の外部依存

| 章 | 外部依存 |
|----|---------|
| ch01-ch04, ch06, ch08, ch11 | なし（標準ライブラリのみ） |
| ch05 | `golang.org/x/sync` |
| ch07 | `github.com/testcontainers/testcontainers-go` |
| ch09 | `go.opentelemetry.io/otel` |
| ch10 | `github.com/sqlc-dev/sqlc`, `github.com/jackc/pgx/v5` |
| ch12 | `google.golang.org/grpc`, `google.golang.org/protobuf` |

## 書籍情報

- **タイトル**: 現場で使えるGo言語実践テクニック
- **サブタイトル**: 設計パターン・並行処理・パフォーマンス最適化
- **著者**: 森川 陽介
- **Amazon**: *近日公開*

## License

MIT License - 詳細は [LICENSE](./LICENSE) を参照してください。
