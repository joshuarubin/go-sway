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

	n, err := client.GetTree(ctx)
	if err != nil {
		return err
	}

	processFocus(ctx, client, n.FocusedNode())

	lifecycle.GoErr(ctx, func() error {
		return sway.Subscribe(ctx, handler{client: client}, sway.EventTypeWindow)
	})

	return lifecycle.Wait(ctx)
}

type handler struct {
	sway.NoOpEventHandler
	client sway.Client
}

func (h handler) Window(ctx context.Context, e sway.WindowEvent) {
	if e.Change != "focus" {
		return
	}

	processFocus(ctx, h.client, e.Container.FocusedNode())
}

func processFocus(ctx context.Context, client sway.Client, node *sway.Node) {
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
