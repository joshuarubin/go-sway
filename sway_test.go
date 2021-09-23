package sway_test

import (
	"context"
	"encoding/json"
	"errors"
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

	n, err := client.GetTree(ctx)
	if err != nil {
		t.Fatal(err)
	}

	processFocus(ctx, client, n.FocusedNode())

	th := testHandler{
		EventHandler: sway.NoOpEventHandler(),
		client:       client,
	}

	if err = sway.Subscribe(ctx, th, sway.EventTypeWindow); err != nil && !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
}

type testHandler struct {
	sway.EventHandler
	client sway.Client
}

func (h testHandler) Window(ctx context.Context, e sway.WindowEvent) {
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
