# FleetDB

> This repository is experimental meaning that it's based on untested ideas or techniques and not yet established or finalized or involves a radically new and innovative style!
> This means that support is best effort (at best!) and we strongly encourage you to NOT use this in production.

FleetDB is a microservice within the Fleet eco-system. It is responsible for providing a store for physical server information. Support to storing the device components that make up the server is available. You are also able to create attributes and versioned-attributes for both servers and the server components.

## Quickstart to running locally

### Install cockroachdb

The cockroachdb client is required to create and drop the test database.

Follow the instructions to install the cockroachdb dependency https://www.cockroachlabs.com/docs/stable/install-cockroachdb.html

You can also run the script:
```bash
./scripts/install_crdb.sh
```
to install version 23.1.14.

### Running fleetdb

To run the fleetdb locally you can bring it up with docker-compose. This will run with released images from the hollow container registry.

```bash
docker-compose -f quickstart.yaml up
```
### Enable tracing

To run the fleetdb locally with tracing enabled you just need to include the `quickstart-tracing.yaml` file.

```bash
docker-compose -f quickstart.yaml -f quickstart-tracing.yaml up
```

### Running with local changes

The `quickstart.yaml` compose file will run fleetdb from released images and not the local code base. If you are doing development and want to run with your local code you can use the following command.

```bash
docker-compose -f quickstart.yaml -f quickstart-dev.yaml up --build
```

NOTE: `--build` is required to get docker-compose to rebuild the container if you have changes. You make also include the `quickstart-tracing.yaml` file if you wish to have tracing support.


### Adding/Changing database schema

Add a new migration file under `db/migrations/` with the schema change

```bash
make docker-up
make test-database
sqlboiler crdb --add-soft-deletes
```

### Run individual integration tests

Export the DB URI required for integration tests.

```bash
export FLEETDB_CRDB_URI="host=localhost port=26257 user=root sslmode=disable dbname=fleetdb_test"
```

Run test.

```bash
go test -timeout 30s -tags testtools -run ^TestIntegrationServerListComponents$ github.com/metal-toolbox/fleetdb/pkg/api/v1 -v
```
