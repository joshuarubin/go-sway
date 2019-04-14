package sway

import "encoding/json"

type Rect struct {
	X      int64 `json:"x,omitempty"`
	Y      int64 `json:"y,omitempty"`
	Width  int64 `json:"width,omitempty"`
	Height int64 `json:"height,omitempty"`
}

type WindowProperties struct {
	Title        string `json:"title,omitempty"`
	Instance     string `json:"instance,omitempty"`
	Class        string `json:"class,omitempty"`
	Role         string `json:"window_role,omitempty"`
	TransientFor int64  `json:"transient_for,omitempty"`
}

type Node struct {
	ID                 int64             `json:"id,omitempty"`
	Name               string            `json:"name,omitempty"`
	Type               string            `json:"type,omitempty"`
	Border             string            `json:"border,omitempty"`
	CurrentBorderWidth int64             `json:"current_border_width,omitempty"`
	Layout             string            `json:"layout,omitempty"`
	Percent            *float64          `json:"percent,omitempty"`
	Rect               Rect              `json:"rect,omitempty"`
	WindowRect         Rect              `json:"window_rect,omitempty"`
	DecoRect           Rect              `json:"deco_rect,omitempty"`
	Geometry           Rect              `json:"geometry,omitempty"`
	Urgent             *bool             `json:"urgent,omitempty"`
	Focused            bool              `json:"focused,omitempty"`
	Focus              []int64           `json:"focus,omitempty"`
	Nodes              []*Node           `json:"nodes,omitempty"`
	FloatingNodes      []*Node           `json:"floating_nodes,omitempty"`
	Representation     *string           `json:"representation,omitempty"`
	AppID              *string           `json:"app_id,omitempty"`
	PID                *uint32           `json:"pid,omitempty"`
	Window             *int64            `json:"window,omitempty"`
	WindowProperties   *WindowProperties `json:"window_properties,omitempty"`
}

func (n *Node) FocusedNode() *Node {
	queue := []*Node{n}
	for len(queue) > 0 {
		n = queue[0]
		queue = queue[1:]

		if n == nil {
			continue
		}

		if n.Focused {
			return n
		}

		queue = append(queue, n.Nodes...)
		queue = append(queue, n.FloatingNodes...)
	}
	return nil
}

type Event interface{}

type WorkspaceEvent struct {
	Change  string `json:"change,omitempty"`
	Current *Node  `json:"current,omitempty"`
	Old     *Node  `json:"old,omitempty"`
}

type WindowEvent struct {
	Change    string `json:"change,omitempty"`
	Container Node   `json:"container,omitempty"`
}

type ShutdownEvent struct {
	Change string `json:"change,omitempty"`
}

type RunCommandReply struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Workspace struct {
	Num     int64  `json:"num,omitempty"`
	Name    string `json:"name,omitempty"`
	Visible bool   `json:"visible,omitempty"`
	Focused bool   `json:"focused,omitempty"`
	Urgent  bool   `json:"urgent,omitempty"`
	Rect    Rect   `json:"rect,omitempty"`
	Output  string `json:"output,omitempty"`
}

type Refresh float64

func (r *Refresh) UnmarshalJSON(raw []byte) error {
	var n int64
	if err := json.Unmarshal(raw, &n); err != nil {
		return err
	}
	*r = Refresh(float64(n) / 1000)
	return nil
}

type OutputMode struct {
	Width   int64   `json:"width,omitempty"`
	Height  int64   `json:"height,omitempty"`
	Refresh Refresh `json:"refresh,omitempty"`
}

type Output struct {
	Name             string       `json:"name,omitempty"`
	Make             string       `json:"make,omitempty"`
	Model            string       `json:"model,omitempty"`
	Serial           string       `json:"serial,omitempty"`
	Active           bool         `json:"active,omitempty"`
	Primary          bool         `json:"primary,omitempty"`
	Scale            float64      `json:"scale,omitempty"`
	Transform        string       `json:"transform,omitempty"`
	CurrentWorkspace string       `json:"current_workspace,omitempty"`
	Modes            []OutputMode `json:"modes,omitempty"`
	CurrentMode      OutputMode   `json:"current_mode,omitempty"`
	Rect             Rect         `json:"rect,omitempty"`
}

type BarConfigGaps struct {
	Top    int64 `json:"top,omitempty"`
	Right  int64 `json:"right,omitempty"`
	Bottom int64 `json:"bottom,omitempty"`
	Left   int64 `json:"left,omitempty"`
}

type BarConfigColors struct {
	Background              string `json:"background,omitempty"`
	Statusline              string `json:"statusline,omitempty"`
	Separator               string `json:"separator,omitempty"`
	FocusedBackground       string `json:"focused_background,omitempty"`
	FocusedStatusline       string `json:"focused_statusline,omitempty"`
	FocusedSeparator        string `json:"focused_separator,omitempty"`
	FocusedWorkspaceText    string `json:"focused_workspace_text,omitempty"`
	FocusedWorkspaceBG      string `json:"focused_workspace_bg,omitempty"`
	FocusedWorkspaceBorder  string `json:"focused_workspace_border,omitempty"`
	ActiveWorkspaceText     string `json:"active_workspace_text,omitempty"`
	ActiveWorkspaceBG       string `json:"active_workspace_bg,omitempty"`
	ActiveWorkspaceBorder   string `json:"active_workspace_border,omitempty"`
	InactiveWorkspaceText   string `json:"inactive_workspace_text,omitempty"`
	InactiveWorkspaceBG     string `json:"inactive_workspace_bg,omitempty"`
	InactiveWorkspaceBorder string `json:"inactive_workspace_border,omitempty"`
	UrgentWorkspaceText     string `json:"urgent_workspace_text,omitempty"`
	UrgentWorkspaceBG       string `json:"urgent_workspace_bg,omitempty"`
	UrgentWorkspaceBorder   string `json:"urgent_workspace_border,omitempty"`
	BindingModeText         string `json:"binding_mode_text,omitempty"`
	BindingModeBG           string `json:"binding_mode_bg,omitempty"`
	BindingModeBorder       string `json:"binding_mode_border,omitempty"`
}

type BarID string

type BarConfigUpdateEvent = BarConfig

type BarConfig struct {
	ID                   BarID           `json:"id,omitempty"`
	Mode                 string          `json:"mode,omitempty"`
	Position             string          `json:"position,omitempty"`
	StatusCommand        string          `json:"status_command,omitempty"`
	Font                 string          `json:"font,omitempty"`
	WorkspaceButtons     bool            `json:"workspace_buttons,omitempty"`
	BindingModeIndicator bool            `json:"binding_mode_indicator,omitempty"`
	Verbose              bool            `json:"verbose,omitempty"`
	Colors               BarConfigColors `json:"colors,omitempty"`
	Gaps                 BarConfigGaps   `json:"gaps,omitempty"`
	BarHeight            int64           `json:"bar_height,omitempty"`
	StatusPadding        int64           `json:"status_padding,omitempty"`
	StatusEdgePadding    int64           `json:"status_edge_padding,omitempty"`
}

type Version struct {
	Major                int64  `json:"major,omitempty"`
	Minor                int64  `json:"minor,omitempty"`
	Patch                int64  `json:"patch,omitempty"`
	HumanReadable        string `json:"human_readable,omitempty"`
	LoadedConfigFileName string `json:"loaded_config_file_name,omitempty"`
}

type Config struct {
	Config string `json:"config,omitempty"`
}

type TickReply struct {
	Success bool `json:"success,omitempty"`
}

type LibInput struct {
	SendEvents      string  `json:"send_events,omitempty"`
	Tap             string  `json:"tap,omitempty"`
	TapButtonMap    string  `json:"tap_button_map,omitempty"`
	TapDrag         string  `json:"tap_drag,omitempty"`
	TapDragLock     string  `json:"tap_drag_lock,omitempty"`
	AccelSpeed      float64 `json:"accel_speed,omitempty"`
	AccelProfile    string  `json:"accel_profile,omitempty"`
	NaturalScroll   string  `json:"natural_scroll,omitempty"`
	LeftHanded      string  `json:"left_handed,omitempty"`
	ClickMethod     string  `json:"click_method,omitempty"`
	MiddleEmulation string  `json:"middle_emulation,omitempty"`
	ScrollMethod    string  `json:"scroll_method,omitempty"`
	ScrollButton    int64   `json:"scroll_button,omitempty"`
	DWT             string  `json:"dwt,omitempty"`
}

type Input struct {
	Identifier          string    `json:"identifier,omitempty"`
	Name                string    `json:"name,omitempty"`
	Vendor              int64     `json:"vendor,omitempty"`
	Product             int64     `json:"product,omitempty"`
	Type                string    `json:"type,omitempty"`
	XKBActiveLayoutName *string   `json:"xkb_active_layout_name,omitempty"`
	LibInput            *LibInput `json:"libinput,omitempty"`
}

type Seat struct {
	Name         string  `json:"name,omitempty"`
	Capabilities int64   `json:"capabilities,omitempty"`
	Focus        int64   `json:"focus,omitempty"`
	Devices      []Input `json:"devices,omitempty"`
}

type ModeEvent struct {
	Change      string `json:"change,omitempty"`
	PangoMarkup bool   `json:"pango_markup,omitempty"`
}

type Binding struct {
	Change         string   `json:"change,omitempty"`
	Command        string   `json:"command,omitempty"`
	EventStateMask []string `json:"event_state_mask,omitempty"`
	InputCode      int64    `json:"input_code,omitempty"`
	Symbol         *string  `json:"symbol,omitempty"`
	InputType      string   `json:"input_type,omitempty"`
}

type BindingEvent struct {
	Change  string  `json:"change,omitempty"`
	Binding Binding `json:"binding,omitempty"`
}

type TickEvent struct {
	First   bool   `json:"first,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type BarStatusUpdateEvent struct {
	ID                string `json:"id,omitempty"`
	VisibleByModifier bool   `json:"visible_by_modifier,omitempty"`
}
