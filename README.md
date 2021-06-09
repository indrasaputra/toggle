# Toggle

[![Go Report Card](https://goreportcard.com/badge/github.com/indrasaputra/toggle)](https://goreportcard.com/report/github.com/indrasaputra/toggle)
[![Workflow](https://github.com/indrasaputra/toggle/workflows/Test/badge.svg)](https://github.com/indrasaputra/toggle/actions)
[![codecov](https://codecov.io/gh/indrasaputra/toggle/branch/main/graph/badge.svg?token=TF36qAeLI0)](https://codecov.io/gh/indrasaputra/toggle)
[![Maintainability](https://api.codeclimate.com/v1/badges/019a5e0793400e5e90ba/maintainability)](https://codeclimate.com/github/indrasaputra/toggle/maintainability)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=indrasaputra_toggle&metric=alert_status)](https://sonarcloud.io/dashboard?id=indrasaputra_toggle)
[![Go Reference](https://pkg.go.dev/badge/github.com/indrasaputra/toggle.svg)](https://pkg.go.dev/github.com/indrasaputra/toggle)

Toggle is a [Feature-Flag](https://martinfowler.com/articles/feature-toggles.html) application.

## Owner

[Indra Saputra](https://github.com/indrasaputra)

## API

### gRPC

The API can be seen in proto files (`*.proto`) in directory [proto](/proto).

### RESTful JSON

The API is automatically generated in OpenAPIv2 format when generating gRPC codes.
The generated files are stored in directory [openapiv2](/openapiv2) in JSON format (`*.json`).
To see the RESTful API contract, do the following:
- Open the generated json file(s)
- Copy the content
- Open [https://editor.swagger.io/](https://editor.swagger.io/)
- Paste the content in [https://editor.swagger.io/](https://editor.swagger.io/)

## How to Run

- Read [Prerequisites](doc/PREREQUISITES.md).
- Then, read [How to Run](doc/HOW_TO_RUN.md).

## Development Guide

- Read [Prerequisites](doc/PREREQUISITES.md).
- Then, read [Development Guide](doc/DEVELOPMENT_GUIDE.md).

## Code Map

- Read [Code Map](doc/CODE_MAP.md)

## Monitoring

The application already emits necessary telemetry. If application's dependencies are run using [docker compose](doc/HOW_TO_RUN.md#docker), then monitoring is [provided by default](docker-compose.yaml). Otherwise, you have to provide them.
These are stacks used as monitoring system.

| Monitoring       | Stack                                      | Address                                           |
| ---              | ---                                        | ---                                               |
| Metrics          | [Prometheus](https://prometheus.io/)       | [http://localhost:9090](http://localhost:9090)    |
| Visualization    | [Grafana](https://grafana.com/)            | [http://localhost:3000](http://localhost:3000)    |
| Tracing          | [Jaeger](https://www.jaegertracing.io/)    | [http://localhost:16686](http://localhost:16686)  |

Special for Grafana, there is [provided dashboard](infrastructure/grafana.dashboard.json) that can be imported. The dashboard contains some basic panels, such as throughput, latency, and error rate.

Currently, tracing only works on gRPC server (handler), service/usecase, and redis. Postgres is not traced yet.