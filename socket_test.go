package sway_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	sway "github.com/joshuarubin/go-sway"
)

func printJSON(v interface{}) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}

func TestSocket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
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

	workspaces, err := client.GetWorkspaces(ctx)
	if err != nil {
		t.Fatal(err)
	}

	printJSON(workspaces)

	outputs, err := client.GetOutputs(ctx)
	if err != nil {
		t.Fatal(err)
	}

	printJSON(outputs)

	marks, err := client.GetMarks(ctx)
	if err != nil {
		t.Fatal(err)
	}

	printJSON(marks)

	barIDs, err := client.GetBarIDs(ctx)
	if err != nil {
		t.Fatal(err)
	}

	printJSON(barIDs)

	for _, id := range barIDs {
		config, err := client.GetBarConfig(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		printJSON(config)
	}

	version, err := client.GetVersion(ctx)
	if err != nil {
		t.Fatal(err)
	}

	printJSON(*version)

	bindingModes, err := client.GetBindingModes(ctx)
	if err != nil {
		t.Fatal(err)
	}

	printJSON(bindingModes)

	config, err := client.GetConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}

	printJSON(config)

	tick, err := client.SendTick(ctx, "foo")
	if err != nil {
		t.Fatal(err)
	}

	printJSON(tick)

	inputs, err := client.GetInputs(ctx)
	if err != nil {
		t.Fatal(err)
	}

	printJSON(inputs)

	seats, err := client.GetSeats(ctx)
	if err != nil {
		t.Fatal(err)
	}

	printJSON(seats)

	fh := focusHandler(client)
	fh(ctx, n.FocusedNode())

	h := sway.EventHandler{
		Window: func(ctx context.Context, e sway.WindowEvent) {
			if e.Change != "focus" {
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
