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
	Class        string `json:"class,omitempty"`
	Instance     string `json:"instance,omitempty"`
	Role         string `json:"window_role,omitempty"`
	Type         string `json:"window_type,omitempty"`
	TransientFor *int64 `json:"transient_for,omitempty"`
}

type Node struct {
	// The internal unique ID for this node
	ID                 int64             `json:"id,omitempty"`

	// The name of the node such as the output name or window title.
	// For the scratchpad, this will be "__i3_scratch" for compatibility with i3.
	Name               string            `json:"name,omitempty"`

	// The node type.
	// It can be "root", "output", "workspace", "con", or "floating_con"
	Type               string            `json:"type,omitempty"`

	// The border style for the node.
	// It can be "normal", "none", "pixel", or "csd"
	Border             string            `json:"border,omitempty"`

	// Number of pixels used for the border width
	CurrentBorderWidth int64             `json:"current_border_width,omitempty"`

	// The node's layout.
	// It can either be "splith", "splitv", "stacked", "tabbed", or "output"
	Layout             string            `json:"layout,omitempty"`

	// The node's orientation.
	// It can be "vertical", "horizontal", or "none"
	Orientation        string            `json:"orientation,omitempty"`

	// The percentage of the node's parent that it takes up or null for the root
	// and other special nodes such as the scratchpad
	Percent            *float64          `json:"percent,omitempty"`

	// The absolute geometry of the node. The window decorations are excluded
	// from this, but borders are included.
	Rect               Rect              `json:"rect,omitempty"`

	// The geometry of the contents inside the node. The window decorations are
	// excluded from this calculation, but borders are included.
	WindowRect         Rect              `json:"window_rect,omitempty"`

	// The geometry of the decorations for the node relative to the parent node
	DecoRect           Rect              `json:"deco_rect,omitempty"`

	// The natural geometry of the contents if it were to size itself
	Geometry           Rect              `json:"geometry,omitempty"`

	// Whether the node or any of its descendants has the urgent hint set.
	// Note: This may not exist when compiled without xwayland support
	Urgent             *bool             `json:"urgent,omitempty"`

	// Whether the node is sticky (shows on all workspaces)
	Sticky             bool              `json:"sticky,omitempty"`

	// List of marks assigned to the node
	Marks              []interface{}     `json:"marks,omitempty"`

	// Whether the node is currently focused by the default seat (seat0)
	Focused            bool              `json:"focused,omitempty"`

	// Array of child node IDs in the current focus order
	Focus              []int64           `json:"focus,omitempty"`

	// The tiling children nodes for the node
	Nodes              []*Node           `json:"nodes,omitempty"`

	// The floating children nodes for the node
	FloatingNodes      []*Node           `json:"floating_nodes,omitempty"`

	// (Only workspaces) A string representation of the layout of the workspace
	// that can be used as an aid in submitting reproduction steps for bug reports
	Representation     *string           `json:"representation,omitempty"`

	// (Only containers and views) The fullscreen mode of the node.
	// 0 means none, 1 means full workspace, and 2 means global fullscreen
	FullscreenMode     *int              `json:"fullscreen_mode,omitempty"`

	// (Only views) For an xdg-shell view, the name of the application, if set.
	// Otherwise, null
	AppID              *string           `json:"app_id,omitempty"`

	// (Only views) The PID of the application that owns the view
	PID                *uint32           `json:"pid,omitempty"`

	// (Only views) Whether the node is visible
	Visible            *bool             `json:"visible,omitempty"`

	// (Only views) The shell of the view, such as "xdg_shell" or "xwayland"
	Shell              *string           `json:"shell,omitempty"`

	// (Only views) Whether the view is inhibiting the idle state
	InhibitIdle        *bool             `json:"inhibit_idle,omitempty"`

	// (Only views) An object containing the state of the application and user
	// idle inhibitors. "application" can be "enabled" or "none".
	// "user" can be "focus", "fullscreen", "open", "visible" or "none".
	IdleInhibitors     interface{}       `json:"idle_inhibitors,omitempty"`

	// (Only xwayland views) The X11 window ID for the xwayland view
	Window             *int64            `json:"window,omitempty"`

	// (Only xwayland views) An object containing the "title", "class", "instance",
	// "window_role", "window_type", and "transient_for" for the view
	WindowProperties   *WindowProperties `json:"window_properties,omitempty"`
}

// FocusedNode traverses the node tree and returns the focused node
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

// WorkspaceEvent is sent whenever a change involving a workspace occurs
type WorkspaceEvent struct {
	// The type of change that occurred
	// The following change types are currently available:
	// init:   the workspace was created
	// empty:  the workspace is empty and is being destroyed since it is not
	//         visible
	// focus:  the workspace was focused. See the old property for the previous
	//         focus
	// move:   the workspace was moved to a different output
	// rename: the workspace was renamed
	// urgent: a view on the workspace has had their urgency hint set or all
	//         urgency hints for views on the workspace have been cleared
	// reload: The configuration file has been reloaded
	Change string `json:"change,omitempty"`

	// An object representing the workspace effected or null for "reload" changes
	Current *Node `json:"current,omitempty"`

	// For a "focus" change, this is will be an object representing the workspace
	// being switched from. Otherwise, it is null
	Old *Node `json:"old,omitempty"`
}

// WindowEvent is sent whenever a change involving a view occurs
type WindowEvent struct {
	// The type of change that occurred
	//
	// The following change types are currently available:
	// new:             The view was created
	// close:           The view was closed
	// focus:           The view was focused
	// title:           The view's title has changed
	// fullscreen_mode: The view's fullscreen mode has changed
	// move:            The view has been reparented in the tree
	// floating:        The view has become floating or is no longer floating
	// urgent:          The view's urgency hint has changed status
	// mark:            A mark has been added or removed from the view
	Change string `json:"change,omitempty"`

	// An object representing the view effected
	Container Node `json:"container,omitempty"`
}

// ShutdownEvent is sent whenever the IPC is shutting down
type ShutdownEvent struct {
	// A string containing the reason for the shutdown.  Currently, the only
	// value for change is "exit", which is issued when sway is exiting.
	Change string `json:"change,omitempty"`
}

type RunCommandReply struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}

type Workspace struct {
	// The workspace number or -1 for workspaces that do not start with a number
	Num     int64  `json:"num,omitempty"`

	// The name of the workspace
	Name    string `json:"name,omitempty"`

	// Whether the workspace is currently visible on any output
	Visible bool   `json:"visible,omitempty"`

	// Whether the workspace is currently focused by the default seat (seat0)
	Focused bool   `json:"focused,omitempty"`

	// Whether a view on the workspace has the urgent flag set
	Urgent  bool   `json:"urgent,omitempty"`

	// The bounds of the workspace. It consists of x, y, width, and height
	Rect    Rect   `json:"rect,omitempty"`

	// The name of the output that the workspace is on
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
	// The name of the output. On DRM, this is the connector
	Name             string       `json:"name,omitempty"`

	// The make of the output
	Make             string       `json:"make,omitempty"`

	// The model of the output
	Model            string       `json:"model,omitempty"`

	// The output's serial number as a hexadecimal string
	Serial           string       `json:"serial,omitempty"`

	// Whether this output is active/enabled
	Active           bool         `json:"active,omitempty"`

	// Whether this output is on/off (via DPMS)
	DPMS             bool         `json:"dpms,omitempty"`

	// For i3 compatibility, this will be false. It does not make sense in Wayland
	Primary          bool         `json:"primary,omitempty"`

	// The scale currently in use on the output or -1 for disabled outputs
	Scale            float64      `json:"scale,omitempty"`

	// The subpixel hinting current in use on the output.
	// This can be "rgb", "bgr", "vrgb", "vbgr", or "none"
	SubpixelHinting  string       `json:"subpixel_hinting,omitempty"`

	// The transform currently in use for the output. This can be "normal", "90",
	// "180", "270", "flipped-90", "flipped-180", or "flipped-270"
	Transform        string       `json:"transform,omitempty"`

	// The workspace currently visible on the output or null for disabled outputs
	CurrentWorkspace string       `json:"current_workspace,omitempty"`

	// An array of supported mode objects.
	// Each object contains "width", "height", and "refresh"
	Modes            []OutputMode `json:"modes,omitempty"`

	// An object representing the current mode containing "width", "height", and "refresh"
	CurrentMode      OutputMode   `json:"current_mode,omitempty"`

	// The bounds for the output consisting of "x", "y", "width", and "height"
	Rect             Rect         `json:"rect,omitempty"`
}

type BarConfigGaps struct {
	Top    int64 `json:"top,omitempty"`
	Right  int64 `json:"right,omitempty"`
	Bottom int64 `json:"bottom,omitempty"`
	Left   int64 `json:"left,omitempty"`
}

// The colors object contains strings which are all #RRGGBBAA representation
// of the color
type BarConfigColors struct {
	// The color to use for the bar background on unfocused outputs
	Background              string `json:"background,omitempty"`

	// The color to use for the status line text on unfocused outputs
	Statusline              string `json:"statusline,omitempty"`

	// 	The color to use for the separator text on unfocused outputs
	Separator               string `json:"separator,omitempty"`

	// The color to use for the background of the bar on the focused output
	FocusedBackground       string `json:"focused_background,omitempty"`

	// The color to use for the status line text on the focused output
	FocusedStatusline       string `json:"focused_statusline,omitempty"`

	// The color to use for the separator text on the focused output
	FocusedSeparator        string `json:"focused_separator,omitempty"`

	// The color to use for the text of the focused workspace button
	FocusedWorkspaceText    string `json:"focused_workspace_text,omitempty"`

	// The color to use for the background of the focused workspace button
	FocusedWorkspaceBG      string `json:"focused_workspace_bg,omitempty"`

	// The color to use for the border of the focused workspace button
	FocusedWorkspaceBorder  string `json:"focused_workspace_border,omitempty"`

	// The color to use for the text of the workspace buttons for the visible
	// workspaces on unfocused outputs
	ActiveWorkspaceText     string `json:"active_workspace_text,omitempty"`

	// The color to use for the background of the workspace buttons for the
	// visible workspaces on unfocused outputs
	ActiveWorkspaceBG       string `json:"active_workspace_bg,omitempty"`

	// The color to use for the border of the workspace buttons for the visible
	// workspaces on unfocused outputs
	ActiveWorkspaceBorder   string `json:"active_workspace_border,omitempty"`

	// The color to use for the text of the workspace buttons for workspaces
	// that are not visible
	InactiveWorkspaceText   string `json:"inactive_workspace_text,omitempty"`

	// The color to use for the background of the workspace buttons for workspaces
	// that are not visible
	InactiveWorkspaceBG     string `json:"inactive_workspace_bg,omitempty"`

	// The color to use for the border of the workspace buttons for workspaces
	// that are not visible
	InactiveWorkspaceBorder string `json:"inactive_workspace_border,omitempty"`

	// The color to use for the text of the workspace buttons for workspaces
	// that contain an urgent view
	UrgentWorkspaceText     string `json:"urgent_workspace_text,omitempty"`

	// The color to use for the background of the workspace buttons for workspaces
	// that contain an urgent view
	UrgentWorkspaceBG       string `json:"urgent_workspace_bg,omitempty"`

	// The color to use for the border of the workspace buttons for workspaces
	// that contain an urgent view
	UrgentWorkspaceBorder   string `json:"urgent_workspace_border,omitempty"`

	// The color to use for the text of the binding mode indicator
	BindingModeText         string `json:"binding_mode_text,omitempty"`

	// The color to use for the background of the binding mode indicator
	BindingModeBG           string `json:"binding_mode_bg,omitempty"`

	// The color to use for the border of the binding mode indicator
	BindingModeBorder       string `json:"binding_mode_border,omitempty"`
}

// BarConfigUpdateEvent is sent whenever a config for a bar changes. The event
// is identical to that of GET_BAR_CONFIG when a bar ID is given as a payload.
type BarConfigUpdateEvent = BarConfig

// Represents the configuration for the bar with the bar ID sent as the payload
type BarConfig struct {
	// The bar ID
	ID                   string          `json:"id,omitempty"`

	// The mode for the bar. It can be "dock", "hide", or "invisible"
	Mode                 string          `json:"mode,omitempty"`

	// The bar's position. It can currently either be "bottom" or "top"
	Position             string          `json:"position,omitempty"`

	// The command which should be run to generate the status line
	StatusCommand        string          `json:"status_command,omitempty"`

	// The font to use for the text on the bar
	Font                 string          `json:"font,omitempty"`

	// Whether to display the workspace buttons on the bar
	WorkspaceButtons     bool            `json:"workspace_buttons,omitempty"`

	// Minimum width in px for the workspace buttons on the bar
	WorkspaceMinWidth    int64           `json:"workspace_min_width,omitempty"`

	// Whether to display the current binding mode on the bar
	BindingModeIndicator bool            `json:"binding_mode_indicator,omitempty"`

	// For i3 compatibility, this will be the boolean value "false".
	Verbose              bool            `json:"verbose,omitempty"`

	// An object containing the #RRGGBBAA colors to use for the bar.
	Colors               BarConfigColors `json:"colors,omitempty"`

	// An object representing the gaps for the bar consisting of "top", "right",
	// "bottom", and "left".
	Gaps                 BarConfigGaps   `json:"gaps,omitempty"`

	// The absolute height to use for the bar or 0 to automatically size based on
	// the font
	BarHeight            int64           `json:"bar_height,omitempty"`

	// The vertical padding to use for the status line
	StatusPadding        int64           `json:"status_padding,omitempty"`

	// The horizontal padding to use for the status line when at the end of an output
	StatusEdgePadding    int64           `json:"status_edge_padding,omitempty"`
}

// Contains version information about the sway process
type Version struct {
	// The major version of the sway process
	Major                int64  `json:"major,omitempty"`

	// The minor version of the sway process
	Minor                int64  `json:"minor,omitempty"`

	// The patch version of the sway process
	Patch                int64  `json:"patch,omitempty"`

	// A human readable version string that will likely contain more useful
	// information such as the git commit short hash and git branch
	HumanReadable        string `json:"human_readable,omitempty"`

	// The path to the loaded config file
	LoadedConfigFileName string `json:"loaded_config_file_name,omitempty"`
}

type Config struct {
	Config string `json:"config,omitempty"`
}

type TickReply struct {
	Success bool `json:"success,omitempty"`
}

// The libinput object describes the device configuration for libinput devices.
// Only properties that are supported for the device will be added to the object.
// In addition to the possible options listed, all string properties may also
// be unknown, in the case that a new option is added to libinput.
// See sway-input(5) for information on the meaning of the possible values.
type LibInput struct {
	// Whether events are being sent by the device.
	// It can be "enabled", "disabled", or "disabled_on_external_mouse"
	SendEvents        string     `json:"send_events,omitempty"`

	// Whether tap to click is enabled. It can be "enabled" or "disabled"
	Tap               string     `json:"tap,omitempty"`

	// The finger to button mapping in use. It can be "lmr" or "lrm"
	TapButtonMap      string     `json:"tap_button_map,omitempty"`

	// Whether tap-and-drag is enabled. It can be "enabled" or "disabled"
	TapDrag           string     `json:"tap_drag,omitempty"`

	// Whether drag-lock is enabled. It can be "enabled" or "disabled"
	TapDragLock       string     `json:"tap_drag_lock,omitempty"`

	// The pointer-acceleration in use
	AccelSpeed        float64    `json:"accel_speed,omitempty"`

	// The acceleration profile in use. It can be "none", "flat", or "adaptive"
	AccelProfile      string     `json:"accel_profile,omitempty"`

	// Whether natural scrolling is enabled. It can be "enabled" or "disabled"
	NaturalScroll     string     `json:"natural_scroll,omitempty"`

	// Whether left-handed mode is enabled. It can be "enabled" or "disabled"
	LeftHanded        string     `json:"left_handed,omitempty"`

	// The click method in use. It can be "none", "button_areas", or "clickfinger"
	ClickMethod       string     `json:"click_method,omitempty"`

	// Whether middle emulation is enabled. It can be "enabled" or "disabled"
	MiddleEmulation   string     `json:"middle_emulation,omitempty"`

	// The scroll method in use.
	// It can be "none", "two_finger", "edge", or "on_button_down"
	ScrollMethod      string     `json:"scroll_method,omitempty"`

	// The scroll button to use when "scroll_method" is "on_button_down".
	// This will be given as an input event code
	ScrollButton      int64      `json:"scroll_button,omitempty"`

	// Whether disable-while-typing is enabled. It can be "enabled" or "disabled"
	DWT               string     `json:"dwt,omitempty"`

	// An array of 6 floats representing the calibration matrix for absolute
	// devices such as touchscreens
	CalibrationMatrix [6]float64 `json:"calibration_matrix,omitempty"`
}

type Input struct {
	// The identifier for the input device
	Identifier           string    `json:"identifier,omitempty"`

	// The human readable name for the device
	Name                 string    `json:"name,omitempty"`

	// The vendor code for the input device
	Vendor               int64     `json:"vendor,omitempty"`

	// The product code for the input device
	Product              int64     `json:"product,omitempty"`

	// The device type. Currently this can be "keyboard", "pointer", "touch",
	// "tablet_tool", "tablet_pad", or "switch"
	Type                 string    `json:"type,omitempty"`

	// (Only keyboards) The name of the active keyboard layout in use
	XKBActiveLayoutName  *string   `json:"xkb_active_layout_name,omitempty"`

	// (Only keyboards) A list a layout names configured for the keyboard
	XKBLayoutNames       []string `json:"xkb_layout_names,omitempty"`

	// (Only keyboards) The index of the active keyboard layout in use
	XKBActiveLayoutIndex *int64    `json:"xkb_active_layout_index,omitempty"`

	// (Only libinput devices) An object describing the current device settings.
	LibInput             *LibInput `json:"libinput,omitempty"`
}

type Seat struct {
	// The unique name for the seat
	Name         string  `json:"name,omitempty"`

	// The number of capabilities that the seat has
	Capabilities int64   `json:"capabilities,omitempty"`

	// The id of the node currently focused by the seat or 0 when the seat is
	// not currently focused by a node (i.e. a surface layer or xwayland
	// unmanaged has focus)
	Focus        int64   `json:"focus,omitempty"`

	// An array of input devices that are attached to the seat.
	// Currently, this is an array of objects that are identical to those
	// returned by GET_INPUTS
	Devices      []Input `json:"devices,omitempty"`
}

// ModeEvent is sent whenever the binding mode changes
type ModeEvent struct {
	// The binding mode that became active
	Change string `json:"change,omitempty"`

	// Whether the mode should be parsed as pango markup
	PangoMarkup bool `json:"pango_markup,omitempty"`
}

type Binding struct {
	// The command associated with the binding
	Command string `json:"command,omitempty"`

	// An array of strings that correspond to each modifier key for the binding
	EventStateMask []string `json:"event_state_mask,omitempty"`

	// For keyboard bindcodes, this is the key code for the binding. For mouse
	// bindings, this is the X11 button number, if there is an equivalent. In
	// all other cases, this will be 0.
	InputCode int64 `json:"input_code,omitempty"`

	// For keyboard bindsyms, this is the bindsym for the binding. Otherwise,
	// this will be null
	Symbol *string `json:"symbol,omitempty"`

	// The input type that triggered the binding. This is either "keyboard" or
	// "mouse"
	InputType string `json:"input_type,omitempty"`
}

// BindingEvent is sent whenever a binding is executed
type BindingEvent struct {
	// Currently this will only be "run"
	Change  string  `json:"change,omitempty"`

	Binding Binding `json:"binding,omitempty"`
}

// TickEvent is sent when first subscribing to tick events or by a SEND_TICK
// message
type TickEvent struct {
	// Whether this event was triggered by subscribing to the tick events
	First bool `json:"first,omitempty"`

	// The payload given with a SEND_TICK message, if any.
	// Otherwise, an empty string
	Payload string `json:"payload,omitempty"`
}

// BarStateUpdateEvent is sent when the visibility of a bar changes due to a
// modifier being pressed
type BarStateUpdateEvent struct {
	// The bar ID effected
	ID string `json:"id,omitempty"`

	// Whether the bar should be made visible due to a modifier being pressed
	VisibleByModifier bool `json:"visible_by_modifier,omitempty"`
}

// Deprecated: BarStatusUpdateEvent is deprecated, use BarStateUpdateEvent instead
type BarStatusUpdateEvent = BarStateUpdateEvent

// InputEvent is sent when something related to the input devices changes.
type InputEvent struct {
	// What has changed
	//
	// The following change types are currently available:
	// added:           The input device became available
	// removed:         The input device is no longer available
	// xkb_keymap:      (Keyboards only) The keymap for the keyboard has changed
	// xkb_layout:      (Keyboards only) The effective layout in the keymap
	//                  has changed
	// libinput_config: (libinput device only) A libinput config option for the
	//                  device changed
	Change string `json:"change,omitempty"`

	// An object representing the input that is identical the ones
	// GET_INPUTS gives
	Input Input `json:"input,omitempty"`
}
