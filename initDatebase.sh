#/bin/bash

docker run -p 6666:3306 -it -d --name=mysql -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7.19
