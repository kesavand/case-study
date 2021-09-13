
FROM golang:1.15-alpine AS builder

RUN apk add build-base

RUN addgroup -S kesavan && adduser -S kesavan -G kesavan

WORKDIR /app
COPY go.mod .
COPY go.sum .
#RUN go mod download
ADD vendor  ./vendor
COPY . .

ARG version=unknown
RUN echo "Setting version: $version"

# Build
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo \
        -o /app/evnt-rcvr \
        -ldflags \
        "-X github.com/iternal/pkg/cli.Version=$version" \
        cmd/evnt-rcvr/main.go

# -------------
# Image creation stage


FROM alpine:3.4

RUN addgroup -S kesavan && adduser -S kesavan -G kesavan
COPY --from=builder /app/evnt-rcvr /app/evnt-rcvr


COPY env /app/.env
EXPOSE 8080
RUN chown -R kesavan:kesavan /app/*
WORKDIR /app
USER 1000
