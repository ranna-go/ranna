# ranna

ランナー - Experimental code runner microservice based on Docker containers.

## ⚠ PLEASE READ BEFORE USE

First of all, this project is currently **work in progress** and not fully finished.  
Also, this service allows **arbitrary code execution in Docker containers**. This will be a high security risk! If you want to use this service, only use this on a separate, encapsulated server environment!

## 🛠 Architecture

Maybe, to make my thought behind the project more clear, here is a little introduction into the project's architecture.

![](https://i.imgur.com/kJyAmso.png)

As you can see, the project is split up in different services.

- **REST API**: The REST API service is the main entrypoint for code execution.
- **Config Provider**: All services need specific configuration. These are obtained by this service.
- **Spec Provider**: ranna works with `specs`, which describe the runner environments for the `Sandbox Provider`. It provides a map of `language` specifiers *(like `go`, or `python3`)* with their specific runner `specs`.
- **Sandbox Provider**: This is the high level API to create a sandbox environment where the passed code can be run inside and the output can be obtained from.
- **Namespace Provider**: This service is responsible for generating unique namespace identifiers which can be used to pass the provided code as file into the sandbox.
- **File Provider**: This service is responsible for creating the nessecary directory structure and the file, containing the code, which is then passed to the sandbox to be executed.

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

The current solution does not use `docker-in-docker`, so, if you use ranna in a Docker container, the specified `RANNA_HOSTROOTDIR` must match the same directory on the host machine. That's why the command above specifies `/var/opt/ranna:/var/opt/ranna` as rootdir volume bind.

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