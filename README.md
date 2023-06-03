# ChatRoom
golangでWebアプリケーションを作ってみる


### dockerイメージを作る
docker build -t chatroom .

### コンテナの起動
docker run --name mysql-container --rm -d -p 3306:3306 chatroom

### コンテナに入る
docker exec -it mysql-container /bin/bash
