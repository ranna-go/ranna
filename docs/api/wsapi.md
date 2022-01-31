# ranna WebSocket API

You can connect to a WebSocket endpoint to the ranna API via the following endpoint.

```
wss://public.ranna.dev/v1/ws
```

The WebSocket API works with the general principle of operations sent by the client side and events sent by the server side. The data is encoded as JSON objects and sent via Text Messages over the WebSocket API.

# Operations

An operation is composed by an `op` code specifying the operation, an optional `nonce` value and operation specific arguments in form of a key-value map. When you pass a `nonce` with an operation, all subsequent events corresponding to this operation will also carry the specified `nonce` value so you are able to group together events.

## Op Codes

| Code | Name   | Description                       |
| ---- | ------ | --------------------------------- |
| `0`  | `PING` | Ping the WebSocket API.           |
| `1`  | `EXEC` | Invoke a code execution.          |
| `2`  | `KILL` | Kill a running execution process. |

## Arguments

### `0` - `PING`

_No arguments._

### `1` - `EXEC`

Please see [models.ExecutionRequest](https://github.com/ranna-go/ranna/blob/master/docs/api/restapi.md#modelsexecutionrequest) to see arguments passed to this operation.

### `2` - `KILL`

| Name    | Type     | Description                        | Required |
| ------- | -------- | ---------------------------------- | -------- |
| `runid` | `string` | The run ID of the running sandbox. | Yes      |

# Events

An event is composed by an event `code` specifying the event type as well as the event `data` payload. When the invoking operation contained a `nonce`, this is also added to the event object.

## Event Codes

| Code | Name    | Description                                           |
| ---- | ------- | ----------------------------------------------------- |
| `0`  | `PONG`  | Response to the `PING` operation.                     |
| `1`  | `ERROR` | An error occurence.                                   |
| `2`  | `SPAWN` | Indicates a spawned execution instance by the client. |
| `3`  | `LOG`   | Indicates a log output from a running execution.      |
| `4`  | `STOP`  | Indicates the finish of an execution.                 |

## Event Data

### `0` - `PONG`

`Pong!`

### `1` - `ERROR`

The error message.

### `2` - `SPAWN`

| Name    | Type     | Description                        |
| ------- | -------- | ---------------------------------- |
| `runid` | `string` | The run ID of the running sandbox. |

### `3` - `LOG`

| Name     | Type      | Description                        |
| -------- | --------- | ---------------------------------- |
| `runid`  | `string`  | The run ID of the running sandbox. |
| `stdout` | `string?` | The `STDOUT` data chunk.           |
| `stderr` | `string?` | The `STDERR` data chunk.           |

### `3` - `STOP`

| Name         | Type     | Description                                      |
| ------------ | -------- | ------------------------------------------------ |
| `runid`      | `string` | The run ID of the running sandbox.               |
| `exectimems` | `int`    | The total time of the execution in milliseconds. |

# Example

Below, you can see a simple example of the message exchange of a code execution.

🔼 00:00

```json
{
  "op": 1,
  "nonce": 123,
  "args": {
    "language": "go",
    "inline_expression": true,
    "code": "import \"time\"\nfor i := 0; i < 5; i++ { time.Sleep(1 * time.Second); println(\"👋\") }"
  }
}
```

🔽 00:00

```json
{
  "op": 1,
  "nonce": 123,
  "args": {
    "language": "go",
    "inline_expression": true,
    "code": "import \"time\"\nfor i := 0; i < 5; i++ { time.Sleep(1 * time.Second); println(\"👋\") }"
  }
}
```

🔽 00:03

```json
{
  "code": 3,
  "nonce": 123,
  "data": {
    "runid": "2a4d3e48e67995e1a7726d79344c96b8ac68ea035b3d25912a50c643da64d4dc",
    "stderr": "👋\n"
  }
}
```

🔽 00:04

```json
{
  "code": 3,
  "nonce": 123,
  "data": {
    "runid": "2a4d3e48e67995e1a7726d79344c96b8ac68ea035b3d25912a50c643da64d4dc",
    "stderr": "👋\n"
  }
}
```

🔽 00:05

```json
{
  "code": 3,
  "nonce": 123,
  "data": {
    "runid": "2a4d3e48e67995e1a7726d79344c96b8ac68ea035b3d25912a50c643da64d4dc",
    "stderr": "👋\n"
  }
}
```

🔽 00:06

```json
{
  "code": 3,
  "nonce": 123,
  "data": {
    "runid": "2a4d3e48e67995e1a7726d79344c96b8ac68ea035b3d25912a50c643da64d4dc",
    "stderr": "👋\n"
  }
}
```

🔽 00:07

```json
{
  "code": 3,
  "nonce": 123,
  "data": {
    "runid": "2a4d3e48e67995e1a7726d79344c96b8ac68ea035b3d25912a50c643da64d4dc",
    "stderr": "👋\n"
  }
}
```

🔽 00:07

```json
{
  "code": 4,
  "nonce": 123,
  "data": {
    "runid": "2a4d3e48e67995e1a7726d79344c96b8ac68ea035b3d25912a50c643da64d4dc",
    "exectimems": 6748
  }
}
```
