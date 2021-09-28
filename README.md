cloud-interview-app
-----

This repo contains a microservice application for use during various interviews at Syndio. Instructions for the interview will be provided seperatly.

### Architecture

A basic gateway api proxies requests to a restful json api serving employee data from Redis and PostgreSQL. The results from the employeesapi are cached in Redis after being retrieved from PostgreSQL.

```
|-----------|      |---------------|
|  GATEWAY  | -->  | EMPLOYEES API | + REDIS/POSTGRES
|-----------|      |---------------|
```

### Repo Layout

```
Basic configuration and instructions:

├── Dockerfile
├── LICENSE
├── README.md
├── compose.yaml

The employees microservice:

├── employees
│   ├── cmd
│   │   └── employeesapi
│   │       └── employeesapi.go   <- the api service
│   └── internal
│       ├── employeesdb
│       │   ├── employeesdb.go    <- the database client
│       │   └── init.sql
│       └── employeeshttp
│           └── employeeshttp.go  <- the http handlers

The gateway microservice:

├── gateway
│   └── cmd
│       └── gateway
│           └── gateway.go        <- the gateway service

Go dependency information:

├── go.mod
└── go.sum
```

### Setup

The project is configured for use with [docker compose](https://docs.docker.com/compose/) so you'll need to have docker running using something like [Docker for Mac](https://docs.docker.com/desktop/mac/install/).

1. Get the repo locally `git clone git@github.com:syndio/cloud-interview-app.git`
2. Run `docker compose up -d --build` to start everything in the background.

You can view request/error logs for all services by running:

`docker compose logs -f gateway employeesapi`

You can view logs for all supporting services by running:

`docker compose logs -f postgres redis`

You can connect to and query the employees database by running:

```
> docker compose exec postgres psql -U dev -d employees
psql (13.4 (Debian 13.4-1.pgdg110+1))
Type "help" for help.

employees=# select * from employees;
 id | title
----+-------
  8 | foo
 10 | bar
(2 rows)
```

You can connect to and query the redis cache by running:

```
> docker compose exec redis redis-cli
127.0.0.1:6379> KEYS *
1) "employees"
127.0.0.1:6379> GET employees
"[{\"id\":8,\"title\":\"foo\"},{\"id\":10,\"title\":\"bar\"}]"
```

### Example Usage

The examples provided here make use of [httpie](https://httpie.io/).

#### Create an Employee

```
> http POST localhost:6540/employees title=engineer
HTTP/1.1 201 Created
Content-Length: 27
Content-Type: application/json
Date: Tue, 28 Sep 2021 15:17:14 GMT

{
    "id": 5,
    "title": "engineer"
}
```

#### List all Employees

```
> http localhost:6540/employees
HTTP/1.1 200 OK
Content-Length: 56
Content-Type: application/json
Date: Tue, 28 Sep 2021 15:23:36 GMT

[
    {
        "id": 5,
        "title": "engineer"
    },
    {
        "id": 6,
        "title": "manager"
    }
]
```

#### Delete an Employee

```
> http DELETE localhost:6540/employees/5
HTTP/1.1 204 No Content
Date: Tue, 28 Sep 2021 15:24:15 GMT
```
