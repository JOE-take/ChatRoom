# ChatRoom
golangでWebアプリケーションを作ってみる


### dockerイメージを作る
docker build -t chatroom .

### コンテナの作成
docker run --name mysql-container --rm -d -p 3306:3306 chatroom

### コンテナの起動
docker start mysql-container

### コンテナに入る
docker exec -it mysql-container /bin/bash
