create database chatroom;
use chatroom;

create table user(
  userName  varchar(256),
  email     varchar(256) primary key,
  password  varchar(256)
);
