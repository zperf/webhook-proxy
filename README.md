# WebHook Proxy

## Usage

Start proxy at `http://127.0.0.1:4000`, assume http proxy is running at `http://127.0.0.1:8899`

```bash
go run main.go
```

It will forward

```
http://127.0.0.1:4000/hooks.foo.bar/abc123
```

To

```
https://hooks.foo.bar/abc123
```

For more flags, type `go run main.go --help` 
