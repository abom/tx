# tx
Accounts transactions challenge (tested with go version of `go version go1.14.2 linux/amd64
`)

## Building

You can build the binary from source or [just install it directly](#installation)

```
go build
```

Then to run the server

```
./tx -path=data/accounts.json
```

The server will be running on `http://localhost:8000` by default.


To see more options:

```
./tx
```

Would output

```
Usage of ./tx:
  -host string
    	server host, defaults to '127.0.0.1' (default "127.0.0.1")
  -path string
    	path to accounts mock file in json format
  -port string
    	server port, defaults to '8000' (default "8000")
  -timeout string
    	transaction timeout e.g. 10s (default "10s")

```

## Installation
To install it directly, you can do

```
go install github.com/abom/tx@latest
```

After this `tx` command should be available (binary would be added to `PATH`).

## API Endpoints

Base URL would be `http://localhost:8000/api/v1`

| Operation                 | URL                       | Method | Data (Request body)                                                                                                                                         |
| ------------------------- | ------------------------- | ------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| List accounts             | `<BASE_URL>/accounts`     | `GET`  |                                                                                                                                                             |
| List certain account      | `<BASE_URL>/account/<id>` | `GET`  |
| Transfer between accounts | `<BAE_URL>/transfer `     | `POST` | JSON of "from", "to" and "amount", e.g. `{ "from": "862fdd01-4d70-4677-93cb-a01fdb0de976", "to": "459f4752-5163-48b3-afff-24b9511eccc2", "amount": "11.1"}` |

## Running tests

Running all tests

```
go test ./...
```

Running API tests:

```
go test github.com/abom/tx/rest
```
