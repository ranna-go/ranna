# ranna

ランナー - Experimental code runner microservice based on Docker containers.

## ⚠ PLEASE READ BEFORE USE

First of all, this project is currently **work in progress** and not fully finished.  
Also, this service allows **arbitrary code execution in Docker containers**. This will be a high security risk! If you want to use this service, only use this on a separate, encapsulated server environment!

## 📃 Todo

👉 Take a look in the [**issue tracker**](https://github.com/ranna-go/ranna/issues).

## 🛠 Architecture

Maybe, to make my thoughts behind the project more clear, here is a little introduction into the project's architecture.

![](https://i.imgur.com/lW0CNPe.png)

As you can see, the project is split up in different services.

- **REST API**: The REST API service is the main entrypoint for code execution.
- **Config Provider**: All services need specific configuration. These are obtained by this service.
- **Spec Provider**: ranna works with `specs`, which describe the runner environments for the `Sandbox Provider`. It provides a map of `language` specifiers *(like `go`, or `python3`)* with their specific runner `specs`.
- **Sandbox Manager**: A higher levbel abstraction to execute code in sandboxes. Also keeps track of running containers to clean them up after teardown.
- **Sandbox Provider**: This is the high level API to create a sandbox environment where the passed code can be run inside and the output can be obtained from.
- **Namespace Provider**: This service is responsible for generating unique namespace identifiers which can be used to pass the provided code as file into the sandbox.
- **File Provider**: This service is responsible for creating the nessecary directory structure and the file, containing the code, which is then passed to the sandbox to be executed.

## 🚀 Setup

👉 Take a look in the [**wiki**](https://github.com/ranna-go/ranna/wiki/%F0%9F%9A%80-Setup).

## 📡 REST API

👉 Take a look in the [**wiki**](https://github.com/ranna-go/ranna/wiki/%F0%9F%93%A1-API).

### 📦 Client Package

ranna also provides a Go client package available in [`pkg/client`](https://github.com/ranna-go/ranna/tree/master/pkg/client).

See the simple [example implementation](https://github.com/ranna-go/ranna/blob/master/examples/client/main.go) how to use the client package.

[Here](https://pkg.go.dev/github.com/ranna-go/ranna#section-directories) you can find some handy documentation for the provided packages.

---

© 2021 Ringo Hoffmann (zekro Development).  
Covered by the MIT License.