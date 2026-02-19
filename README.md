# go-http-server

Go で HTTP サーバーを実装するサンプルリポジトリ。

## 起動方法

```sh
docker compose up
```

## 動作確認コマンド

```sh
curl -X POST localhost:8000/core/v1/task/get -d '{ "id": "id_01" }'
curl -X POST localhost:8000/core/v1/task/list
curl -X POST localhost:8000/core/v1/task/create -d '{ "title": "title_01" }'
curl -X POST localhost:8000/core/v1/task/update -d '{ "id": "id_01", "title": "title_updated", "status": "DONE" }'
curl -X POST localhost:8000/core/v1/task/delete -d '{ "id": "id_01" }'
```

## ドキュメント

- [API 設計](docs/api-design.md)
- [アーキテクチャ](docs/architecture.md)
- [開発フロー](docs/dev-flow.md)
