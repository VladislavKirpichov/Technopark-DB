FROM golang:1.19 AS build

ADD . /internal
WORKDIR /internal

COPY . .
RUN go mod tidy
RUN go build  ./cmd/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "create user forum_user with superuser password 'ac322f35-e71e-47dd-a580-894bf3e6c460';" &&\
    createdb -O forum_user forum &&\
    /etc/init.d/postgresql stop

EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
USER root

WORKDIR /usr/src/internal
COPY . .
COPY --from=build /internal/main/ .

EXPOSE 5000
ENV PGPASSWORD ac322f35-e71e-47dd-a580-894bf3e6c460
CMD service postgresql start && psql -h localhost -d forum -U forum_user -p 5432 -a -q -f ./db/db.sql && ./main
