package sway

import (
	"context"
)

// EventType is used to choose which events to Subscribe to
type EventType string

const (
	// EventTypeWorkspace is sent whenever an event involving a workspace occurs
	// such as initialization of a new workspace or a different workspace gains
	// focus
	EventTypeWorkspace EventType = "workspace"

	// EventTypeMode is sent whenever the binding mode changes
	EventTypeMode EventType = "mode"

	// EventTypeWindow is sent whenever an event involving a view occurs such as
	// being reparented, focused, or closed
	EventTypeWindow EventType = "window"

	// EventTypeBarConfigUpdate is sent whenever a bar config changes
	EventTypeBarConfigUpdate EventType = "barconfig_update"

	// EventTypeBinding is sent when a configured binding is executed
	EventTypeBinding EventType = "binding"

	// EventTypeShutdown is sent when the ipc shuts down because sway is exiting
	EventTypeShutdown EventType = "shutdown"

	// EventTypeTick is sent when an ipc client sends a SEND_TICK message
	EventTypeTick EventType = "tick"

	// EventTypeBarStateUpdate is sent when the visibility of a bar should change
	// due to a modifier
	EventTypeBarStateUpdate EventType = "bar_state_update"

	// Deprecated: EventTypeBarStatusUpdate is deprecated
	// you should use EventTypeBarStateUpdate instead
	EventTypeBarStatusUpdate EventType = EventTypeBarStateUpdate

	// EventTypeInput is sent when something related to input devices changes
	EventTypeInput EventType = "input"
)

// An EventHandler is passed to Subscribe and its methods are called in response
// to sway events
type EventHandler interface {
	Workspace(context.Context, WorkspaceEvent)
	Mode(context.Context, ModeEvent)
	Window(context.Context, WindowEvent)
	BarConfigUpdate(context.Context, BarConfigUpdateEvent)
	Binding(context.Context, BindingEvent)
	Shutdown(context.Context, ShutdownEvent)
	Tick(context.Context, TickEvent)
	BarStateUpdate(context.Context, BarStateUpdateEvent)
	BarStatusUpdate(context.Context, BarStatusUpdateEvent)
	Input(context.Context, InputEvent)
}

// NoOpEventHandler is used to help provide empty methods that aren't intended
// to be handled by Subscribe
//
//	type handler struct {
//		sway.EventHandler
//	}
//
//	func (h handler) Window(ctx context.Context, e sway.WindowEvent) {
//		...
//	}
//
//	func main() {
//		h := handler{
//			EventHandler: sway.NoOpEventHandler(),
//		}
//
//		ctx := context.Background()
//
//		sway.Subscribe(ctx, h, sway.EventTypeWindow)
//	}
func NoOpEventHandler() EventHandler {
	return noOpEventHandler{}
}

type noOpEventHandler struct{}

func (h noOpEventHandler) Workspace(context.Context, WorkspaceEvent)             {}
func (h noOpEventHandler) Mode(context.Context, ModeEvent)                       {}
func (h noOpEventHandler) Window(context.Context, WindowEvent)                   {}
func (h noOpEventHandler) BarConfigUpdate(context.Context, BarConfigUpdateEvent) {}
func (h noOpEventHandler) Binding(context.Context, BindingEvent)                 {}
func (h noOpEventHandler) Shutdown(context.Context, ShutdownEvent)               {}
func (h noOpEventHandler) Tick(context.Context, TickEvent)                       {}
func (h noOpEventHandler) BarStateUpdate(context.Context, BarStateUpdateEvent) {}
func (h noOpEventHandler) BarStatusUpdate(context.Context, BarStatusUpdateEvent) {}
func (h noOpEventHandler) Input(context.Context, InputEvent) {}

// Subscribe the IPC connection to the events listed in the payload
func Subscribe(ctx context.Context, handler EventHandler, events ...EventType) error {
	n, err := New(ctx)
	if err != nil {
		return err
	}

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
	case eventTypeBarStateUpdate:
		var e BarStateUpdateEvent
		if err := msg.Decode(&e); err == nil {
			h.BarStateUpdate(ctx, e)
		}
	case eventTypeInput:
		var e InputEvent
		if err := msg.Decode(&e); err == nil {
			h.Input(ctx, e)
		}
	}
}
