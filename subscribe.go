package sway

import (
	"context"
)

type EventType string

const (
	EventTypeWorkspace       EventType = "workspace"
	EventTypeMode            EventType = "mode"
	EventTypeWindow          EventType = "window"
	EventTypeBarConfigUpdate EventType = "barconfig_update"
	EventTypeBinding         EventType = "binding"
	EventTypeShutdown        EventType = "shutdown"
	EventTypeTick            EventType = "tick"
	EventTypeBarStatusUpdate EventType = "bar_status_update"
)

type EventHandler struct {
	Workspace func(context.Context, WorkspaceEvent)
	Window    func(context.Context, WindowEvent)
	Shutdown  func(context.Context, ShutdownEvent)
}

func Subscribe(ctx context.Context, handler EventHandler, events ...EventType) error {
	n, err := New(ctx)
	if err != nil {
		return err
	}
	defer n.Close()

	c := n.(*client)

	if err = c.subscribe(ctx, events...); err != nil {
		return err
	}

	for {
		reply, err := c.recvMsg(ctx)
		if err != nil {
			return err
		}

		handler.process(ctx, reply)
	}
}

func (h EventHandler) process(ctx context.Context, reply *message) {
	switch reply.Type {
	case eventReplyTypeWorkspace:
		if h.Workspace == nil {
			return
		}

		var e WorkspaceEvent
		if err := reply.Decode(&e); err != nil {
			return
		}

		h.Workspace(ctx, e)
	case eventReplyTypeMode:
	case eventReplyTypeWindow:
		if h.Window == nil {
			return
		}

		var e WindowEvent
		if err := reply.Decode(&e); err != nil {
			return
		}

		h.Window(ctx, e)
	case eventReplyTypeBarConfigUpdate:
	case eventReplyTypeBinding:
	case eventReplyTypeShutdown:
		if h.Shutdown == nil {
			return
		}

		var e ShutdownEvent
		if err := reply.Decode(&e); err != nil {
			return
		}

		h.Shutdown(ctx, e)
	case eventReplyTypeTick:
	case eventReplyTypeBarStatusUpdate:
	}
}
