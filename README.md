## gracefulhttp
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/aoliveti/gracefulhttp)
[![codecov](https://codecov.io/gh/aoliveti/gracefulhttp/graph/badge.svg?token=j9a2QoWNA5)](https://codecov.io/gh/aoliveti/gracefulhttp)
[![Go Reference](https://pkg.go.dev/badge/github.com/aoliveti/curling)](https://pkg.go.dev/github.com/aoliveti/gracefulhttp)
[![Go Report Card](https://goreportcard.com/badge/github.com/aoliveti/gracefulhttp)](https://goreportcard.com/report/github.com/aoliveti/gracefulhttp)
![GitHub License](https://img.shields.io/github/license/aoliveti/gracefulhttp)

This package extends the functionality of the standard `http.Server` in Go to include best practice configurations and graceful shutdown.

- Smoothly stop accepting new connections while processing existing requests to completion.
- Set a customizable timeout for graceful shutdown, with a default of 5 seconds.
- Forcibly close ongoing connections after the timeout expires.
- Support graceful shutdown for HTTP and HTTPS servers.
- **Apply best practice configurations** inspired by Cloudflare for timeout management and TLS setup ([Read more](https://blog.cloudflare.com/exposing-go-on-the-internet/)).

## Installation
```bash
go get github.com/aoliveti/gracefulhttp
