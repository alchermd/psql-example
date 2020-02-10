## Installation

1. Create a PostgreSQL database named `example`

```bash
$ createdb example
```

2. Update the connection string in the `main()` function on [main.go](https://github.com/alchermd/psql-example/blob/master/main.go)

3. Run or build the app:

```bash
$ go run main.go # or...
$ go build main.go
$ ./main
$ curl -i http://127.0.0.1:8000/
HTTP/1.1 200 OK
Date: Mon, 10 Feb 2020 09:20:08 GMT
Content-Length: 265
Content-Type: text/html; charset=utf-8
...
```

Released under [MIT](LICENSE)