# ws-hello

An example of a Windows Service Application template

## Build

```
> env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/ws-hello.exe .
```
