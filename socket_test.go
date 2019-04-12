package sway_test

import (
	"context"
	"log"
	"testing"
	"time"

	sway "github.com/joshuarubin/go-sway"
)

func TestSocket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := sway.New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	n, err := client.GetTree(ctx)
	if err != nil {
		t.Fatal(err)
	}

	fh := focusHandler(client)

	fh(ctx, n.FocusedNode())

	h := sway.EventHandler{
		Window: func(ctx context.Context, e sway.WindowEvent) {
			if e.Change != sway.WindowChangeFocus {
				return
			}
			fh(ctx, e.Container.FocusedNode())
		},
	}

	err = sway.Subscribe(ctx, h, sway.EventTypeWindow, sway.EventTypeShutdown)
	if err != context.DeadlineExceeded && err != nil {
		t.Fatal(err)
	}
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
