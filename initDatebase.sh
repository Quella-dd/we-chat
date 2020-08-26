#/bin/bash

eval $(docker-machine env default)

docker run -p 6666:3306 -it -d --name=mysql1 -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7.19
