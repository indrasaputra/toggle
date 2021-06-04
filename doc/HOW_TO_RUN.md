## How to Run

There are two ways to run the application. The first is to setup all dependencies manually.
The second is to use docker image and docker compose.

### Manual

- Create `.env` file

    You can copy the `env.example` and change its values accordingly

    ```
    $ cp env.example .env
    ```

- Fill all env variables with prefix `POSTGRES_` according to your PostgreSQL settings

- Run or start PostgreSQL

- Fill the `REDIS_ADDRESS` and `REDIS_TTL`

    `REDIS_TTL` is a Time to Live for a key-value in redis. `REDIS_TTL=5` means its TTL is 5 minutes

- Run or start Redis

- Fill `PORT_GRPC` and `PORT_REST` value as you wish. We use `8080` as default value for `PORT_GRPC` and `8081` for `PORT_REST`.
    `PORT_GRPC` is a port for HTTP/2 gRPC. `PORT_REST` is port for HTTP/1.1 REST.
    We encourage to let both values as default

- Download the dependencies

    ```
    $ make tidy
    ```

- It is always good to have your database migration up-to-date.
    Run the following command to make your database stays up-to-date with the current migrations.

    ```
    $ make migrate url=<postgres url>
    ```

    e.g:

    ```
    $ make migrate url="postgres://user:password@host:port/dbname"
    ```

    **DON'T FORGET** to supply/change the `user`, `password`, `host`, `port`, and `dbname` according to your database settings.

- Run the application

    ```
    $ go run cmd/server/main.go
    ```

### Docker

- Install [Docker Compose](https://docs.docker.com/compose/).

- Run docker compose

    ```
    $ docker-compose up
    ```

- In another terminal window, download the dependencies for the application

    ```
    $ make tidy
    ```

- Create `.env` file

    You can copy the `env.example` and change its values accordingly

    ```
    $ cp env.example .env
    ```

- Run the application

    ```
    $ go run cmd/server/main.go
    ```