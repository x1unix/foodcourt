# voracity mysql db service, version 1.0.0
FROM mysql:5.6
COPY ./my.cnf /etc/mysql/conf.d/
MAINTAINER Denys Sedchenko <dsedchenko@llnw.com>

ENV MYSQL_DATABASE=voracity \
MYSQL_ROOT_PASSWORD=root

ADD ./data.sql /docker-entrypoint-initdb.d
# ADD ./user.sql /docker-entrypoint-initdb.d
EXPOSE 3306
