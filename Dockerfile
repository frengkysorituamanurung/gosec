ARG ALPINE_VERSION=3.19

# stage 1: build go app
FROM golang:1.21-alpine${ALPINE_VERSION} AS builder

# golang build env
ARG CGO_ENABLED=0
ARG GOOS=linux
ENV CGO_ENABLED=${CGO_ENABLED}
ENV GOOS=${GOOS}

WORKDIR /build-stage

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o backend server.go

# stage 2: final image
FROM alpine:${ALPINE_VERSION}

WORKDIR /app

COPY --from=builder /build-stage/backend .

CMD [ "./backend" ]
