## Code Map

Code Map explains how the codes are structured in this project. This document will list down all folders or packages and their purpose.

---

### `bin`

This folder contains any executable binary to support the project.
For example, there is `generate-mock.sh`. It is a shell script file used to generate mocks for all interfaces available in this project.

---

### `cmd`

This folder contains the `main.go`.
The use case may be run and served in multi forms, such as API, cron, or fullstack web.
To cater that case, `cmd` folder can contains subfolders with each folder named accordingly to the form and contain only main package.
e.g: `cmd/api/main.go`, `cmd/cron/main.go`, and `cmd/web/main.go`

For this project, we prefer to use `cmd/server/main.go` as our use cases are only in the form of gRPC server.

---

### `db/migrations`

This folder contains all database migration files. Each migration has exactly two files: UP and DOWN.

---

### `doc`

This folder contains all documents related to the project.

---

### `entity`

This folder contains the domain of the project.
Mostly, this folder contains only structs, constants, global variables, enumerations, or functions with simple logic related to the core domain of the module (not a business logic).
Since we use Protocol Buffer, entity has a close (tightly coupled) relationship with any struct generated from `.proto` files.

---

### `feature`

This folder contains all [Cucumber](https://cucumber.io/docs/guides/) definitions for the purpose of integration / API test.
We use [Gherkin syntax](https://cucumber.io/docs/gherkin/) to define the features.

---

### `internal`

All APIs/codes in the internal folder (and all if its subfolders) are designed to [not be able to be imported](https://golang.org/doc/go1.4#internalpackages).
This folder contains all detail implementation specified in the `service` folder.

---

### `internal/builder`

This folder contains the [builder design pattern](https://sourcemaking.com/design_patterns/builder).
It composes all codes needed to build a full usecase.

---

### `internal/config`

This folder contains configuration for the project.

---

### `internal/grpc/handler`

This folder contains the HTTP/2 gRPC handlers.
Codes in this folder implement gRPC server interface.

---

### `internal/repository`

This folder contains codes that connect to the repository, such as database.
Repository is not limited to databases. Anything that can be a repository can be put here.

---

### `internal/server`

This folder contains all codes needed to define a gRPC and its REST Gateway server.

---

### `openapiv2`

This folder contains API definition for HTTP/1.1 REST.
The contents of this folder are generated from `.proto` as well.

---

### `proto`

This folder contains `.proto` files and all files generated from or based on `.proto`.

---

### `service`

This folder contains the main business logic of the project. Almost all interfaces and all the business logic flows are defined here.
If someone wants to know the flow of the project, they better start to open this folder.

---

### `test`

This folder contains test related stuffs.
For the case of unit test, the unit test files are put in the same directory as the files that they test. It is one of the Go best practice, so we follow.

---

### `test/fixture`

This folder contains a well defined support for test.

---

### `test/mock`

This folder contains mock for testing.

---