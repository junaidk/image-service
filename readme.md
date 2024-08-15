# Running Application

Run `make image` to build app image.

Run `make up` to start app and its dependencies (postgress container).

Im some cases migrat container will fail to run if DB is not ready.
In taht case

Run `make restart` to restart the stack.

Run `make logs` to view app logs.

Run `sample/test.sh` to execute the APIs with sample input.

# API

See `openapi.yml` file
