# Personal Property Prototype

[![Build status](https://img.shields.io/circleci/project/github/transcom/mymove/master.svg)](https://circleci.com/gh/transcom/mymove/tree/master)

This repository contains the application source code for the Personal Property Prototype, a possible next generation version of the Defense Personal Property System (DPS). DPS is an online system managed by the U.S. [Department of Defense](https://www.defense.gov/) (DoD) [Transportation Command](http://www.ustranscom.mil/) (USTRANSCOM) and is used by service members and their families to manage household goods moves.

This prototype was built by a [Defense Digital Service](https://www.dds.mil/) team in support of USTRANSCOM's mission.

## Development

### Project location

All of Go's tooling expects Go code to be checked out in a specific location. Please read about [Go workspaces](https://golang.org/doc/code.html#Workspaces) for a full explanation. If you just want to get started, then decide where you want all your go code to live and configure the GOPATH environment variable accordingly. For example, if you want your go code to live at `~/code/go`, you should add the following like to your `.bash_profile`:

```bash
export GOPATH=~/code/go
```

If you are OK with using the default location for go code (`~/go`), then there is nothing to do. Since this is the default location, using it means you do not need to set `$GOPATH` yourself.

*Regardless of where your go code is located*, you need to add `$GOPATH/bin` to your `PATH` so that executables installed with the go tooling can be found. Add the following to your `.bash_profile`:

```bash
export PATH=$(go env GOPATH)/bin:$PATH
```

Once that's done, you have go installed, and you've re-sourced your profile, you can checkout this repository by running `go get github.com/transcom/mymove/cmd/webserver` (This will emit an error "can't load package:" but will have cloned the source correctly). You will then find the code at `$GOPATH/src/github.com/transcom/mymove`

If you have already checked out the code somewhere else, you can just move it to be in the above location and everything will work correctly.

### Project Layout

All of our code is intermingled in the top level directory of mymove. Here is an explanation of what some of these directories contain:

`bin`: A location for tools helpful for developing this project \
`build`: The build output directory for the client. This is what the development server serves \
`cmd`: The location of main packages for any go binaries we build (right now, just webserver) \
`config`: Config files can be dropped here \
`docs`: A location for docs for the project. This is where ADRs are \
`migrations`: Database migrations live here \
`node_modules`: Cached dependencies for the client \
`pkg`: The location of all of our go libraries, most of our go code lives here \
`public`: The client's static resources \
`src`: The react source code for the client \
`vendor`: Cached dependencies for the server

### Prerequisites

* Install Go with Homebrew. Make sure you do not have other installations.
* Run `bin/prereqs` and install everything it tells you to. *Do not configure postgres to automatically start at boot time!*
* Run `make deps`.
* [EditorConfig](http://editorconfig.org/) allows us to manage editor configuration (like indent sizes,) with a [file](https://github.com/transcom/ppp/blob/master/.editorconfig) in the repo. Install the appropriate plugin in your editor to take advantage of that.

### Setup: Database

You will need to setup a local database before you can begin working on the local server / client. Docker will need to be running for any of this to work.

1. `make db_dev_init`: initializes a Docker container with a Postgres database.
1. `make db_dev_migrate`: runs all existing database migrations, which do things like creating table structures, etc.
1. You can validate that your dev database is running by running `bin/psql-dev`. This puts you in a postgres shell. Type `\dt` to show all tables, and `\q` to quit.

### Setup: Server

1. `make server_run`: installs dependencies and builds both the client and the server, then runs the server.

For faster development, use `make server_run_dev`. This builds and runs the server but skips updating dependences and re-building the client. Those tasks can be accomplished as needed with `make server_deps` and `make client_build`

You can verify the server is working as follows:

`> curl http://localhost:8080/api/v1/issues --data "{ \"body\": \"This is a test issue\"}"`

from which the response should be like

`{"id":"d5735bc0-7553-4d80-a42d-ea1e50bbcfc4", "body": "This is a test issue", "created_at": "2018-01-04 14:47:28.894988", "updated_at": "2018-01-04 14:47:28.894988"}`

Dependencies are managed by [dep](https://github.com/golang/dep). New dependencies are automatically detected in import statements. To add a new dependency to the project, import it in a source file and then run `dep ensure`

### Setup: Client

1. `make server_run`
1. `make client_run`

The above will start the server running and starts the webpack dev server, proxied to our running go server.

Dependencies are managed by yarn. To add a new dependency, use `yarn add`

### API

The api is defined in a single file: ./swagger.yaml and served at /api/v1/swagger.yaml. it is the single source of truth for what the API contract between client and server should be.

### Testing

There are a few handy targets in the Makefile to help you run tests:

* `make client_test`: Run frontend testing suites.
* `make server_test`: Run backend testing suites.
* `make test`: Run both client- and server-side testing suites.

### Database

#### Dev Commands

There are a few handy targets in the Makefile to help you interact with the dev database:

* `make db_dev_init`: Initializes a new postgres Docker container with a test database and runs it. You must do this before any other database operations.
* `make db_dev_run`: Starts the previously initialized Docker container if it has been stopped.
* `make db_dev_reset`: Destroys your database container. Useful if you want to start from scratch.
* `make db_dev_migrate`: Applies database migrations against your running database container.

#### Migrations

If you need to change the database schema, you'll need to write a migration.

Creating a migration:

Use soda (a part of [pop](https://github.com/markbates/pop/)) to generate migrations. In order to make using soda easy, a wrapper is in `./bin/soda` that sets the go environment and working directory correctly.

If you are generating a new model, use `./bin/soda generate model model-name column-name:type column-name:type ...`. id, created_at and updated_at are all created automatically.

If you are modifying an existing model, use `./bin/soda generate migration migration-name` and add the [Fizz instructions](https://github.com/markbates/pop/blob/master/fizz/README.md) yourself to the created files.

Running migrations in local development:

1. Use `make db_dev_migrate` to run migrations against your local dev environment.

Running migrations on Staging / Production:

Migrations are run automatically by CircleCI as part of the standard deploy process.

1. CircleCI builds and registers a container that includes the `soda` binary, along with migrations files.
1. CircleCI deploys this container to ECS and runs it as a one-off 'task'.
1. Migrations run inside the container against the environment's database.
1. If migrations fail, CircleCI fails the deploy.
1. If migrations pass, CircleCI continues with the deploy.

### Troubleshooting

* Random problems may arise if you have old Docker containers running. Run `docker ps` and if you see containers unrelated to our app, consider stopping them.
* If you have problems connecting to postgres, or running related scripts, make sure you aren't already running a postgres daemon. You can check this by typing `ps aux | grep postgres` and looking for existing processes.
* If you happen to have installed pre-commit in a virtual environment not with brew, running bin/prereqs will not alert you. You may run into issues when running `make deps`. To install pre-commit: `brew install pre-commit`.
