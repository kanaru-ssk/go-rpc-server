# アーキテクチャ

## 概要

[The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)をベースに設計

<svg width="256" height="256" viewBox="0 0 256 256" xmlns="http://www.w3.org/2000/svg">
    <circle cx="50%" cy="50%" r="50%"   fill="#8FBAE5" />
    <circle cx="50%" cy="50%" r="37.5%" fill="#8FE5A1" />
    <circle cx="50%" cy="50%" r="25%"   fill="#E58F8F" />
    <circle cx="50%" cy="50%" r="12.5%" fill="#EAE5A0" />
    <text x="50%" y="86"  text-anchor="middle" dominant-baseline="middle">usecase</text>
    <text x="50%" y="54"  text-anchor="middle" dominant-baseline="middle">interface</text>
    <text x="50%" y="50%" text-anchor="middle" dominant-baseline="middle">entity</text>
</svg>

依存方向を同心円の内向きのみにする。

```
/
├── cmd/httpserver/main.go
├── entity/
├── interface/
│   ├── inbound/
│   └── outbound/
├── lib/
└── usecase/
```

entity はビジネス固有のロジック、usecase はアプリケーション固有の処理フローを、外部の DB や API、フレームワークなどに依存させずに実装する。

interface/inbound は、ユーザーのリクエストをパースし、usecase を呼び出した結果を json や html などの適切な形に整形してユーザーに返す。

外部の DB や API リクエストが必要なメソッドは entity で interface のみを定義しておき、interface/outbound で実装する。

## entity

ビジネスルールをカプセル化する。

DB や外部 API の操作が必要なメソッドは interface のみ定義し、特定の DB や外部 API に依存させない。

usecase を見れば全体の処理の流れを追える状態にするため、極力複数のエンティティに跨るドメインサービスは使用しない。

## usecase

エンティティ を使用してアプリケーション固有の処理フローを記述する。

## interface

### inbound

外部からのリクエストを受け取り、usecase を呼び出し、結果に応じてレスポンスを返す。

### outbound

entity 層で定義した interface の実装を行う。DB や外部 API などとのやり取りを実装する。

## 各レイヤーの依存方向

entity, usecase が API の通信形式や DB などに依存しないようにする。
