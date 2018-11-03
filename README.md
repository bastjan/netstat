# netstat for go

[![GoDoc][doc-img]][doc] [![CI Status][ci-img]][ci] [![Coverage Status][cover-img]][cover] [![Go Report Card][report-img]][report]

Package netstat helps you query open network connections.

## Getting Started

```go
import "github.com/bastjan/netstat"

// Query open tcp sockets
netstat.TCP.Connections()

// Query open udp sockets for ipv6 connections
netstat.UDP6.Connections()
```

See [netstat_tulpen.go](examples/netstat_tulpen/netstat_tulpen.go) for a more throughout example.

## Development Status: Stable

This library is v1 and follows SemVer.

No breaking changes will be made to exported APIs before v2.0.0.

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
