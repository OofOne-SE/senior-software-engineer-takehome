# Weather data API (WD-API)

Hey people!

How did I accomplish this task? Image the following process.
- I talked to the frontend guys and we agreed on a swagger openapi specification for this go service. Take a look at [OpenAPI specification](./spec/openapi.yaml)
- I generated a [client](./client), [service](./restapi) and a [model](./models).
- I wired [the endpoint implementation](./internal/weather/weather.go) in [a service implementation](./restapi/configure_backend.go) and created [some simple integration tests using generated client](./internal/weather/weather_test.go)
- Afterwards I implemented [the database call backs](./internal/postgres/postgres.go) for postgres db.

## How to test the code?
I assume you are working with some linux machine and installed docker, docker compose, golang...
- Spin up the env with `docker compose build && docker compose up -d`
- Start the websocket listener with `python connect_to_websocket.py`
- If you would like to import the given data run `python import_data.py` or `make test` to start the integration tests

Have fun!