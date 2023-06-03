FROM mysql:8.0

COPY db/chatroom.sql /docker-entrypoint-initdb.d/chatroom.sql

ENV MYSQL_ROOT_PASSWORD=password

