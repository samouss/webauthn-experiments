# Webauthn Experiments

Webauthn Experiments built with Go, CRA, Bulma.

Live: [https://webauthn-experiments.herokuapp.com/](https://webauthn-experiments.herokuapp.com/)


## Installation

Clone the repository and then run the following command:

```
cd webauthn-website
yarn
```

## Run the development application

For run the development server in watch mode on `localhost:8080`:

```
cd webauthn-api
go run .
```

Then in a other tab you can run the development server for the client in watch mode on `localhost:3000`:

```
cd webauthn-website
yarn start
```

## Run the production application

You can build the production application with Docker:

```
docker build -t webauthn-api .
```

Some environement variables might be required to run the production version:

```
HOST=0.0.0.0
ORIGIN=https://example.com
SESSION_KEY=RANDOM_VALUE
SESSION_SECURE=true
STATIC_DIR=webauthn-website
```

You can take a look at the [`main.go`](webauthn-api/main.go) for more information.
