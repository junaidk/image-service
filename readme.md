# Image Service

HTTP REST API for uploading images via an invitation link.

User can gererate expireable links and can use the links to upload images.

It also provide some statistics about the uploaded images.


## Running Application

Provided docker compose file starts the db and run the migration command.

And after that start the application process.

Run `make image` to build app image.

Run `make up` to start app and its dependencies (postgress container).

Im some cases migrat container will fail to run if DB is not ready.
In taht case

Run `make restart` to restart the stack.

Run `make logs` to view app logs.

Run `sample/test.sh` to execute the APIs with sample input.

# API

See `openapi.yml` file

# Project Structure

The project organizes code with the following approach:

- Application domain types go in the root—User.
- Implementations of the application domain go in subpackages — `postgress`, `http`, etc.
- Everything is tied together in the cmd subpackages — `cmd/server` .
