# tx
Accounts transactions challenge (tested with go version of `go version go1.14.2 linux/amd64
`)

## Building

```
go build
```

Then to run the server

```
./tx -path=/path/to/accounts-mock.json
```

To see more options:

```
./tx
```

To install

```
go install github.com/abom/tx
```

The server will be running on `http://localhost:8000`

## API Endpoints

Base URL would be `http://localhost:8000/api/v1`

| Operation                 | URL                       | Method | Data (Request body)                                                                                                                                         |
| ------------------------- | ------------------------- | ------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| List accounts             | `<BASE_URL>/accounts`     | `GET`  |                                                                                                                                                             |
| List certain account      | `<BASE_URL>/account/<id>` | `GET`  |
| Transfer between accounts | `<BAE_URL>/transfer `     | `POST` | JSON of "from", "to" and "amount", e.g. `{ "from": "862fdd01-4d70-4677-93cb-a01fdb0de976", "to": "459f4752-5163-48b3-afff-24b9511eccc2", "amount": "11.1"}` |

## Running tests

Only processor tests for now, to run them

```
go install github.com/abom/tx
go test github.com/abom/tx/processor
```
