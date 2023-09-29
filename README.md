# Project Managment Demo Backend

## How to start

Using docker-compose:

```
docker-compose up --build
```

## API

### REST API

Each widget has own REST API in the following format:

#### backend.url.com/api/{widget}/routes

See more details for each widget api in:

- api/kanban.go
- api/gantt.go
- api/todo.go
- api/scheduler.go

### WS API

Each widget has own Web Socket API in the following format:

#### backend.url.com/api/{widget}/v1

See more details for each widget in:

- publisher/kanban.go
- publisher/gantt.go
- publisher/todo.go
- publisher/scheduler.go
