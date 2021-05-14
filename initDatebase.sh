#!/bin/bash

# 创建测试用的mysql, redis容器环境
docker run -p 6666:3306 -it -d --name=we_chat -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7.19

docker run -p 6379:6379 -it -d --name=we_chat_redis redis --requirepass "123456"
