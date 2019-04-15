[![Go Report Card](https://goreportcard.com/badge/github.com/joshuarubin/go-sway)](https://goreportcard.com/report/github.com/joshuarubin/go-sway) [![GoDoc](https://godoc.org/github.com/joshuarubin/go-sway?status.svg)](https://godoc.org/github.com/joshuarubin/go-sway)

This package simplifies working with the [sway](https://swaywm.org/) IPC from Go.
It was highly influenced by the [i3 package](https://github.com/i3/go-i3).

While the i3 and sway IPCs share much in common, they are not identical. This package provides the complete sway api.

## Differences from the i3 package

* Retries are not handled. Use tools like systemd to automatically restart apps that use this library.
* A much simpler interface for subscriptions and handling events.
* No global state.
* Use of Context throughout.

## Assumptions

* The `$SWAYSOCK` variable must be set properly in the environment
* sway is running on a machine with the same byteorder as the client
