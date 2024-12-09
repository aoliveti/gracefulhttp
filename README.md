## gracefulhttp
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/aoliveti/gracefulhttp)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/aoliveti/gracefulhttp/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/aoliveti/curling)](https://pkg.go.dev/github.com/aoliveti/gracefulhttp)
[![codecov](https://codecov.io/gh/aoliveti/gracefulhttp/graph/badge.svg?token=j9a2QoWNA5)](https://codecov.io/gh/aoliveti/gracefulhttp)
[![Go Report Card](https://goreportcard.com/badge/github.com/aoliveti/gracefulhttp)](https://goreportcard.com/report/github.com/aoliveti/gracefulhttp)
![GitHub License](https://img.shields.io/github/license/aoliveti/gracefulhttp)

This package extends the functionality of the standard [http.Server](https://pkg.go.dev/net/http#Server) in Go to include best practice configurations and graceful shutdown.

- Smoothly stop accepting new connections while processing existing requests to completion.
- Set a customizable timeout for graceful shutdown, with a default of 5 seconds.
- Forcibly close ongoing connections after the timeout expires.
- Support graceful shutdown for HTTP and HTTPS servers.
- **Apply best practice configurations** inspired by Cloudflare for timeout management and TLS setup ([Read more](https://blog.cloudflare.com/exposing-go-on-the-internet/)).

## Installation
```bash
go get github.com/aoliveti/gracefulhttp
```

## Usage

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"

	"github.com/aoliveti/gracefulhttp"
	"golang.org/x/sys/unix"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGINT, unix.SIGTERM)
	defer cancel()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello world!"))
	})

	srv := gracefulhttp.Bind(":8080", nil)

	log.Println("starting server...")

	if err := srv.ListenAndServeWithShutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("graceful shutdown completed successfully")
}
```
You can instantiate a GracefulServer in two different ways:
```go
gracefulhttp.GracefulServer{
    Server: http.Server{
        Addr:    ":8080",
        Handler: nil,
    },
}
```
or by using the Bind() function:
```go
func Bind(addr string, handler http.Handler) *GracefulServer
```

## Listen and serve
To start the HTTP server with the provided address and handler, you need to use this function:
```go
func (s *GracefulServer) ListenAndServeWithShutdown(ctx context.Context, opts ...GracefulServerOption) error
```
If instead you need to start a server that accepts connections over HTTPS, you must ensure to provide files containing a certificate and matching private key for the server. You should use this function:
```go
func (s *GracefulServer) ListenAndServeTLSWithShutdown(ctx context.Context, certFile string, keyFile string, opts ...GracefulServerOption) error
```
It's possible to pass options to set timeouts and [TLS configuration](https://pkg.go.dev/crypto/tls#Config). Here's a summary table:

| Option                  | Description                                                                                                       |
|-------------------------|-------------------------------------------------------------------------------------------------------------------|
| WithShutdownTimer       | Sets the timeout for a graceful shutdown, after which all active connections will be forcibly closed              |
| WithCloudflareTimeouts  | Applies timeout patches to the server, implementing best practice configurations inspired by Cloudflare           |
| WithCloudflareTLSConfig | Applies TLS configuration patches to the server, implementing best practice configurations inspired by Cloudflare |
| WithTLSConfig           | Sets the provided TLS configuration                                                                               |

## Build and Test

### Building the Project

You can build the project using the provided [Makefile](Makefile). Simply run:

```sh
make
```

### Running Tests

To run tests, use the following command:

```sh
make test
```

Before running the tests, a certificate and key are generated using OpenSSL and placed in the `certs` directory.
The generated certificate and key are used for testing HTTPS functions.

For detailed documentation on OpenSSL, visit https://www.openssl.org/

## License

The library is released under the MIT license. See [LICENSE](LICENSE) file.
