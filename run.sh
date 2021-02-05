#!/bin/bash

mysqlroot=root

docker run --name=mysqls -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=thodo --network thodo-net --network-alias mysql-net -d -p 3306:3306 mysql/mysql-server:latest
sleep 2 # let time for mysql to be completely setup
docker exec mysqls mysql --enable-cleartext-plugin -p$mysqlroot <<< "UPDATE mysql.user SET host = '%' WHERE user = 'root'"
docker build -t thodo .
docker run --network thodo-net -t thodo:latest