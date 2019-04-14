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

type EventHandler interface {
	Workspace(context.Context, WorkspaceEvent)
	Mode(context.Context, ModeEvent)
	Window(context.Context, WindowEvent)
	BarConfigUpdate(context.Context, BarConfigUpdateEvent)
	Binding(context.Context, BindingEvent)
	Shutdown(context.Context, ShutdownEvent)
	Tick(context.Context, TickEvent)
	BarStatusUpdate(context.Context, BarStatusUpdateEvent)
}

type NoOpEventHandler struct{}

func (h NoOpEventHandler) Workspace(context.Context, WorkspaceEvent)             {}
func (h NoOpEventHandler) Mode(context.Context, ModeEvent)                       {}
func (h NoOpEventHandler) Window(context.Context, WindowEvent)                   {}
func (h NoOpEventHandler) BarConfigUpdate(context.Context, BarConfigUpdateEvent) {}
func (h NoOpEventHandler) Binding(context.Context, BindingEvent)                 {}
func (h NoOpEventHandler) Shutdown(context.Context, ShutdownEvent)               {}
func (h NoOpEventHandler) Tick(context.Context, TickEvent)                       {}
func (h NoOpEventHandler) BarStatusUpdate(context.Context, BarStatusUpdateEvent) {}

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
		msg, err := c.recvMsg(ctx)
		if err != nil {
			return err
		}

		processEvent(ctx, handler, msg)
	}
}

func processEvent(ctx context.Context, h EventHandler, msg *message) {
	switch msg.Type {
	case eventTypeWorkspace:
		var e WorkspaceEvent
		if err := msg.Decode(&e); err == nil {
			h.Workspace(ctx, e)
		}
	case eventTypeMode:
		var e ModeEvent
		if err := msg.Decode(&e); err == nil {
			h.Mode(ctx, e)
		}
	case eventTypeWindow:
		var e WindowEvent
		if err := msg.Decode(&e); err == nil {
			h.Window(ctx, e)
		}
	case eventTypeBarConfigUpdate:
		var e BarConfigUpdateEvent
		if err := msg.Decode(&e); err == nil {
			h.BarConfigUpdate(ctx, e)
		}
	case eventTypeBinding:
		var e BindingEvent
		if err := msg.Decode(&e); err == nil {
			h.Binding(ctx, e)
		}
	case eventTypeShutdown:
		var e ShutdownEvent
		if err := msg.Decode(&e); err == nil {
			h.Shutdown(ctx, e)
		}
	case eventTypeTick:
		var e TickEvent
		if err := msg.Decode(&e); err == nil {
			h.Tick(ctx, e)
		}
	case eventTypeBarStatusUpdate:
		var e BarStatusUpdateEvent
		if err := msg.Decode(&e); err == nil {
			h.BarStatusUpdate(ctx, e)
		}
	}
}
