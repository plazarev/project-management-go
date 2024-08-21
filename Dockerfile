FROM golang:1.19-bullseye AS builder
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /project-manager-go -a -ldflags '-linkmode external -extldflags "-static"'

FROM scratch
WORKDIR /app
COPY --from=builder /project-manager-go /app/project-manager-go
COPY ./demodata /app/demodata
COPY ./config.yml /app/config.yml

EXPOSE 3000

CMD ["/app/project-manager-go"]
