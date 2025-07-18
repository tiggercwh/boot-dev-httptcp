# Multiple protocols demo in Go

This repository contains a collection of Go modules to demonstrate the fundamentals of the HTTP protocol, TCP, and UDP.

## Components

This project includes the following modules:

- **`httpserver`**: A simple HTTP server that listens on port `2025`. It includes handlers for various routes to demonstrate different HTTP features.
- **`tcplistener`**: A raw TCP listener on port `2025`. It accepts connections, parses incoming data as HTTP requests, and prints the request details to the console. This is useful for inspecting raw HTTP requests.
- **`udpsender`**: A basic UDP client that sends data from standard input to `localhost:2025`.

## Getting Started

You just need to have Go installed on your machine to run this repo.

## Usage Examples

Here are a few examples of how to interact with the applications.

### Interacting with `httpserver`

Run the `httpserver` in one terminal:

```bash
go run ./cmd/httpserver
```

In another terminal, you can use `curl` to send requests to it.

**Example: Get a successful response**

```bash
curl -v http://localhost:2025/normalreq
```

**Example: Get an error response**

```bash
curl -v http://localhost:2025/myproblem
```

**Example: Use the httpbin proxy**

This will proxy your request to `https://httpbin.org/get`.

```bash
curl -v http://localhost:2025/httpbin/get
```

### Visualizing `tcplistener`

Run the `tcplistener` in one terminal:

```bash
go run ./cmd/tcplistener
```

In another terminal, send an HTTP request using `curl`. The `-v` flag shows the request being sent.

```bash
curl -v -X POST -d "hello world" http://localhost:2025/foo
```

The `tcplistener` terminal will print the raw request it received, for example:

```
Accepted connection from 127.0.0.1:54622
...
Body:
hello world
Connection to  127.0.0.1:54622 closed
```

### Visualizing `udpsender`

To see the output of the `udpsender`, you can use `netcat` (`nc`) to act as a UDP server.

In one terminal, start `netcat` to listen for UDP packets on port `2025`:

```bash
nc -ul -p 2025
```

In another terminal, run the `udpsender`:

```bash
go run ./cmd/udpsender
```

Now, type any message into the `udpsender` terminal and press Enter. You will see the message appear in the `netcat` terminal.

### Build

You can build all the applications at once using the Go toolchain:

```bash
go build ./...
```

This will create the executables in your current directory.
