
FROM golang:1.22.0-alpine AS build

WORKDIR /WorkerQueueSystem

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /WorkerQueueSystem/WorkerQueueSystemExec ./main.go

FROM alpine:latest

WORKDIR /WorkerQueueSystem

COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/

COPY --from=build /WorkerQueueSystem/WorkerQueueSystemExec /WorkerQueueSystem/

EXPOSE 8080

ENTRYPOINT ["/WorkerQueueSystem/WorkerQueueSystemExec"]
