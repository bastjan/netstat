# netstat for go

[![GoDoc][doc-img]][doc] [![CI Status][ci-img]][ci] [![Coverage Status][cover-img]][cover] [![Go Report Card][report-img]][report]

Package netstat helps you query open socket connections.

## Getting Started

```go
import "github.com/bastjan/netstat"

// Query open tcp sockets
netstat.TCP.Entries()

// Query open udp sockets for ipv6 connections
netstat.UDP6.Entries()
```

## Development Status: Work in Progress

The api is not yet final and can change.
First stable release will be version 1.0.0.

## Support for Mac OS and *BSD

There is currently no support planned for MacOS or *BSD without procfs support.

[doc]: https://godoc.org/github.com/bastjan/netstat
[doc-img]: https://godoc.org/github.com/bastjan/netstat?status.svg
[cover]: https://codecov.io/gh/bastjan/netstat
[cover-img]: https://codecov.io/gh/bastjan/netstat/branch/master/graph/badge.svg
[ci]: https://travis-ci.org/bastjan/netstat
[ci-img]: https://travis-ci.org/bastjan/netstat.svg?branch=master
[report]: https://goreportcard.com/report/github.com/bastjan/netstat
[report-img]: https://goreportcard.com/badge/github.com/bastjan/netstat
