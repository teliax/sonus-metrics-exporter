FROM golang:1.18-alpine as build
LABEL maintainer "Teliax"

RUN apk --no-cache add ca-certificates \
     && apk --no-cache add --virtual build-deps git gcc musl-dev

COPY ./ /go/src/github.com/teliax/sonus-metrics-exporter
WORKDIR /go/src/github.com/teliax/sonus-metrics-exporter

RUN go get \
 && go test ./... \
 && go build -o /bin/main

FROM alpine:3

RUN apk --no-cache add ca-certificates \
     && addgroup exporter \
     && adduser -S -G exporter exporter
USER exporter
COPY --from=build /bin/main /bin/main
ENV LISTEN_PORT=9172
EXPOSE 9172
ENTRYPOINT [ "/bin/main" ]
