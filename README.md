# ws-hello

An example of a Windows Service Application template

## Build

```
> env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/ws-hello.exe .
```

## Create Service for Windows
```
> sc create [nameOfService] [binpath=pathToExecFile]
[type={own|share|kernel|filesys|rec|adapt|interact type={own|share}}]
[start={boot|system|auto|demand|disabled}] [displayname=nombreDescriptivo]
```

Example
```
> sc.exe create WsHello binpath="D:\Documents\Download\ws-hello-master\build\ws-hello.exe" displayName="Window Service Hello Application" start=auto
```

## Delete Service

```
> sc delete WsHello 
```
