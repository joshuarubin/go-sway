package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall"

	sway "github.com/joshuarubin/go-sway"
	"github.com/joshuarubin/lifecycle"
)

func main() {
	if err := run(); err != nil && !isSignal(err) {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func isSignal(err error, sigs ...os.Signal) bool {
	serr, ok := err.(lifecycle.ErrSignal)
	if !ok {
		return false
	}
	switch serr.Signal {
	case syscall.SIGINT, syscall.SIGTERM:
		return true
	}
	return false
}

func run() error {
	ctx := lifecycle.New(context.Background())

	client, err := sway.New(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	fh := focusHandler(client)

	n, err := client.GetTree(ctx)
	if err != nil {
		return err
	}

	fh(ctx, n.FocusedNode())

	h := sway.EventHandler{
		Window: func(ctx context.Context, e sway.WindowEvent) {
			if e.Change != "focus" {
				return
			}
			fh(ctx, e.Container.FocusedNode())
		},
	}

	lifecycle.GoErr(ctx, func() error {
		return sway.Subscribe(ctx, h, sway.EventTypeWindow)
	})

	return lifecycle.Wait(ctx)
}

func focusHandler(client sway.Client) func(context.Context, *sway.Node) {
	return func(ctx context.Context, node *sway.Node) {
		if node == nil {
			return
		}

		opt := "none"
		if node.AppID == nil || *node.AppID != "kitty" {
			opt = "altwin:ctrl_win"
		}

		if _, err := client.RunCommand(ctx, `input '*' xkb_options `+opt); err != nil {
			log.Println(err)
		}
	}
}
