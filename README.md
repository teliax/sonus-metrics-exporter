# Prometheus Sonus Exporter

Exposes basic metrics for your Sonus SBC from the API, to a Prometheus
compatible endpoint.  Based on the
[Github Exporter](https://github.com/infinityworks/github-exporter) code.

## Configuration

This exporter is setup to take input from environment variables:

### Required
* `API_URL` The URL of the Sonus API.  Should appear as `https://{ip}/api`
* `API_USER` The username to use when logging in to the Sonus API.
* `API_PASSWORD` The password to use when authenticating to the Sonus.

### Optional
* `LISTEN_PORT` The port you wish to run the container on, the Dockerfile defaults this to `9172`
* `METRICS_PATH` the metrics URL path you wish to use, defaults to `/metrics`
* `LOG_LEVEL` The level of logging the exporter will run with, defaults to `debug`


## Install and deploy

Run manually from Docker Hub:
```
docker run -d --restart=always -p 9172:9172 -e API_USER="username" -e API_PASSWORD="password" teliax/sonus-metrics-exporter
```

Build a docker image:
```
docker build -t <image-name> .
docker run -d --restart=always -p 9172:9172  -e API_USER="username" -e API_PASSWORD="password" <image-name>
```

## Docker compose

```
sonus-metrics-exporter:
    tty: true
    stdin_open: true
    expose:
      - 9172
    ports:
      - 9172:9172
    image: teliax/sonus-metrics-exporter:latest
    environment:
      - API_USER=username
      - API_PASSWORD=password

```

## Metrics

Metrics will be made available on port 9172 by default
An example of these metrics can be found in the `METRICS.md` markdown file in the root of this repository

## Metadata
