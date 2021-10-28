# Toggle

[![Go Report Card](https://goreportcard.com/badge/github.com/indrasaputra/toggle)](https://goreportcard.com/report/github.com/indrasaputra/toggle)
[![Workflow](https://github.com/indrasaputra/toggle/workflows/Test/badge.svg)](https://github.com/indrasaputra/toggle/actions)
[![codecov](https://codecov.io/gh/indrasaputra/toggle/branch/main/graph/badge.svg?token=TF36qAeLI0)](https://codecov.io/gh/indrasaputra/toggle)
[![Maintainability](https://api.codeclimate.com/v1/badges/019a5e0793400e5e90ba/maintainability)](https://codeclimate.com/github/indrasaputra/toggle/maintainability)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=indrasaputra_toggle&metric=alert_status)](https://sonarcloud.io/dashboard?id=indrasaputra_toggle)
[![Go Reference](https://pkg.go.dev/badge/github.com/indrasaputra/toggle.svg)](https://pkg.go.dev/github.com/indrasaputra/toggle)

Toggle is a [Feature-Flag](https://martinfowler.com/articles/feature-toggles.html) application. It uses event-driven paradigm.

## Owner

[Indra Saputra](https://github.com/indrasaputra)

## API

### gRPC

The API can be seen in proto files (`*.proto`) in directory [proto](proto/indrasaputra/toggle/v1/toggle.proto).

### RESTful JSON

The API is automatically generated in OpenAPIv2 format when generating gRPC codes.
The generated files are stored in directory [openapiv2](openapiv2) in JSON format (`*.json`).
To see the RESTful API contract, do the following:
- Open the generated json file(s), such as [toggle.swagger.json](openapiv2/proto/indrasaputra/toggle/v1/toggle.swagger.json)
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

## Testing

### Unit Test

```
$ make test.unit
```

### Integration Test

[godog](https://github.com/cucumber/godog/#install) is mandatory to perform integration test.

To run the integration test, make sure you already run the application successfully. Follow [How to Run](doc/HOW_TO_RUN.md) for the guideline.
When application is running, then run command to execute integration test.

```
$ make test.integration
```

You can also set the server URL, in case your default server is not localhost.

```
$ SERVER_URL=http://toggle:8081/v1/toggles make test.integration
```

### Load Test

Running smoke test, load test, and stress test is encouraged to know the sanity, performance, and stability of the service.
[k6](https://k6.io/docs/) is used as load test executor.

```sh
$ k6 run <path to script file>
```

e.g:

```sh
$ k6 run internal/script/loadtest/load_test.js
```

or use docker

```sh
$ make test.load
```

## Observability

The application already emits necessary telemetry. If application's dependencies are run using [docker compose](doc/HOW_TO_RUN.md#docker), then monitoring is [provided by default](docker-compose.yaml). Otherwise, you have to provide them.
These are stacks used as monitoring system.

| Monitoring       | Stack                                      | Address                                           |
| ---              | ---                                        | ---                                               |
| Metrics          | [Prometheus](https://prometheus.io/)       | [http://localhost:9090](http://localhost:9090)    |
| Visualization    | [Grafana](https://grafana.com/)            | [http://localhost:3000](http://localhost:3000)    |
| Tracing          | [Jaeger](https://www.jaegertracing.io/)    | [http://localhost:16686](http://localhost:16686)  |
| Log              | [Zap](https://github.com/uber-go/zap)      | Stdout                                            |

Currently, tracing only works on gRPC server (handler), service/usecase, and redis. Postgres is not traced yet.

## SDK

There is already an SDK to access Toggle. Currently, the SDK only supports Go. The SDK codes are located in this very repository. Visit [Toggle SDK](pkg/sdk/toggle/client.go).
