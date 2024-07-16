FROM golang:1.21 AS build

ARG POLIGONO_VERSION=main

WORKDIR /tmp/vega-core

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${POLIGONO_VERSION}" -o ./bin/vega-core ./server.go

FROM busybox:1.35.0-uclibc as busybox

FROM gcr.io/distroless/base-debian11 AS build-release-stage

ENV PORT=8888

ENV AUTHENTICATION_ENABLE_RBAC=true
ENV AUTHENTICATION_TYPE=BasicAuth
ENV ENABLE_AUTHENTICATION=true
ENV TRINO_DSN=http://user@trino:8080?catalog=default&schema=test
ENV GIN_MODE=debug

COPY --from=busybox /bin/sh /bin/sh

COPY --from=build /tmp/vega-core/bin/vega-core /app/vega-core

EXPOSE ${PORT}

CMD ["/app/vega-core"]