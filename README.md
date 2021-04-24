# ranna

ランナー - Experimental code runner microservice based on Docker containers.

## ⚠ PLEASE READ BEFORE USE

First of all, this project is currently **work in progress** and not fully finished.  
Also, this service allows **arbitrary code execution in Docker containers**. This will be a high security risk! If you want to use this service, only use this on a separate, encapsulated server environment!

## 🚀 Setup

First of all, you need to know that ranna needs access to the docker socket. 
> This can be done over network, but currently, there is no network file provider implementation to push the source code files to the docker host system.

You can get the binary directly by compiling the source code.
```
$ go build -o ranna cmd/ranna/main.go
```

Or build tzhe provided Docker image.
```
$ docker build . -t ranna
```

But keep in mind, when you are using the Docker image, you need to pass the host's docker socket as volume to the container.
```
$ docker run --name ranna \
    -e RANNA_HOSTROOTDIR=/var/opt/ranna \
    -v /var/opt/ranna:/var/opt/ranna \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -p 8080:8080
    ranna
```

## 📡 REST API

### `GET /v1/spec`

Returns a map of runner environment specifications where the key is
the `language` specifier and the value is the `spec`.

```
> GET /v1/spec HTTP/2
```

The response of this request will look like following:

```
< HTTP/2 200 OK
< Server: ranna
< Content-Type: application/json
< Content-Length: 155
```
```json
{
  "go": {
    "image": "golang:alpine",
    "entrypoint": "go run",
    "filename": "main.go"
  },
  "python3": {
    "image": "python:alpine",
    "entrypoint": "python3",
    "filename": "main.py"
  }
}
```

### `POST /v1/exec`

Execute code.

```
> POST /v1/exec HTTP/2
> Content-Type: application/json
```
```json
{
  "language": "python3",
  "code": "print('Hello world!')",
  "arguments": ["some", "crazy", "arguments"],
  "environment": {
    "MYVAR": "my value"
  }
}
```

The response of this request will look like following:

```
< HTTP/2 200 OK
< Server: ranna
< Content-Type: application/json
< Content-Length: 39
```
```json
{
  "stdout": "Hello world!\n",
  "stderr": ""
}
```

### 📦 Client Package

ranna also provides a Go client package available in [`pkg/client`](https://github.com/zekroTJA/ranna/tree/master/pkg/client).

See the simple [example implementation](https://github.com/zekroTJA/ranna/blob/master/examples/client/main.go) how to use the client package.

[Here](https://pkg.go.dev/github.com/zekroTJA/ranna#section-directories) you can find some handy documentation for the provided packages.

---

© 2021 Ringo Hoffmann (zekro Development).  
Covered by the MIT License.