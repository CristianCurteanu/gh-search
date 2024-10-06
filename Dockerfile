FROM golang:1.22.4-alpine3.19 AS builder
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/runner main.go

FROM alpine:3.19
COPY --chown=65534:65534 --from=builder /go/bin/runner .
USER 65534

ENTRYPOINT [ "./runner" ]