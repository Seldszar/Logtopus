# Logtopus

> Server for receiving and storing remote application logs

## Requirements

- [Go](https://go.dev/dl)

## Usage

```sh
$ go run main.go
```

### Send log messages

Send new messages to `POST /{name}`, where `{name}` is your application name with the following JSON body:

| Key       | Type                                   |
|-----------|----------------------------------------|
| `level`   | "DEBUG" \| "INFO" \| "WARN" \| "ERROR" |
| `message` | string                                 |
| `detail`  | object                                 |

### Access log messages

Received messages are stored in a `logs` folder, where the filename is the application name.

## License

Copyright (c) Alexandre Breteau

This software is released under the terms of the MIT License.
See the [LICENSE](LICENSE) file for further information.
