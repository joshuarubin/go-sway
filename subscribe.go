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
	Workspace       func(context.Context, WorkspaceEvent)
	Mode            func(context.Context, ModeEvent)
	Window          func(context.Context, WindowEvent)
	BarConfigUpdate func(context.Context, BarConfigUpdateEvent)
	Binding         func(context.Context, BindingEvent)
	Shutdown        func(context.Context, ShutdownEvent)
	Tick            func(context.Context, TickEvent)
	BarStatusUpdate func(context.Context, BarStatusUpdateEvent)
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
	case eventTypeWorkspace:
		if h.Workspace == nil {
			return
		}

		var e WorkspaceEvent
		if err := reply.Decode(&e); err == nil {
			h.Workspace(ctx, e)
		}
	case eventTypeMode:
		if h.Mode == nil {
			return
		}

		var e ModeEvent
		if err := reply.Decode(&e); err == nil {
			h.Mode(ctx, e)
		}
	case eventTypeWindow:
		if h.Window == nil {
			return
		}

		var e WindowEvent
		if err := reply.Decode(&e); err == nil {
			h.Window(ctx, e)
		}
	case eventTypeBarConfigUpdate:
		if h.BarConfigUpdate == nil {
			return
		}

		var e BarConfigUpdateEvent
		if err := reply.Decode(&e); err == nil {
			h.BarConfigUpdate(ctx, e)
		}
	case eventTypeBinding:
		if h.Binding == nil {
			return
		}

		var e BindingEvent
		if err := reply.Decode(&e); err == nil {
			h.Binding(ctx, e)
		}
	case eventTypeShutdown:
		if h.Shutdown == nil {
			return
		}

		var e ShutdownEvent
		if err := reply.Decode(&e); err == nil {
			h.Shutdown(ctx, e)
		}
	case eventTypeTick:
		if h.Tick == nil {
			return
		}

		var e TickEvent
		if err := reply.Decode(&e); err == nil {
			h.Tick(ctx, e)
		}
	case eventTypeBarStatusUpdate:
		if h.BarStatusUpdate == nil {
			return
		}

		var e BarStatusUpdateEvent
		if err := reply.Decode(&e); err == nil {
			h.BarStatusUpdate(ctx, e)
		}
	}
}
