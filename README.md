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
