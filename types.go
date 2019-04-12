package sway

type NodeID int64

type NodeType string

const (
	NodeTypeRoot        NodeType = "root"
	NodeTypeOutput      NodeType = "output"
	NodeTypeCon         NodeType = "con"
	NodeTypeFloatingCon NodeType = "floating_con"
	NodeTypeWorkspace   NodeType = "workspace"
)

type BorderStyle string

const (
	BorderStyleNormal BorderStyle = "normal"
	BorderStyleNone   BorderStyle = "none"
	BorderStylePixel  BorderStyle = "pixel"
	BorderStyleCSD    BorderStyle = "csd"
)

type Layout string

const (
	LayoutSplitH  Layout = "splith"
	LayoutSplitV  Layout = "splitv"
	LayoutStacked Layout = "stacked"
	LayoutTabbed  Layout = "tabbed"
	LayoutOutput  Layout = "output"
)

type Rect struct {
	X      int64 `json:"x"`
	Y      int64 `json:"y"`
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

type WindowProperties struct {
	Title        string `json:"title"`
	Instance     string `json:"instance"`
	Class        string `json:"class"`
	Role         string `json:"window_role"`
	TransientFor NodeID `json:"transient_for"`
}

type Node struct {
	ID                 NodeID            `json:"id"`
	Name               string            `json:"name"`
	Type               NodeType          `json:"type"`
	Border             BorderStyle       `json:"border"`
	CurrentBorderWidth int64             `json:"current_border_width"`
	Layout             Layout            `json:"layout"`
	Percent            *float64          `json:"percent"`
	Rect               Rect              `json:"rect"`
	WindowRect         Rect              `json:"window_rect"`
	DecoRect           Rect              `json:"deco_rect"`
	Geometry           Rect              `json:"geometry"`
	Urgent             *bool             `json:"urgent"`
	Focused            bool              `json:"focused"`
	Focus              []NodeID          `json:"focus"`
	Nodes              []*Node           `json:"nodes"`
	FloatingNodes      []*Node           `json:"floating_nodes"`
	Representation     *string           `json:"representation"`
	AppID              *string           `json:"app_id"`
	PID                *uint32           `json:"pid"`
	Window             *int64            `json:"window"`
	WindowProperties   *WindowProperties `json:"window_properties"`
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

type WorkspaceChange string

const (
	WorkspaceChangeInit   WorkspaceChange = "init"
	WorkspaceChangeEmpty  WorkspaceChange = "empty"
	WorkspaceChangeFocus  WorkspaceChange = "focus"
	WorkspaceChangeMove   WorkspaceChange = "move"
	WorkspaceChangeRename WorkspaceChange = "rename"
	WorkspaceChangeUrgent WorkspaceChange = "urgent"
	WorkspaceChangeReload WorkspaceChange = "reload"
)

type WorkspaceEvent struct {
	Change  WorkspaceChange `json:"change"`
	Current *Node           `json:"current"`
	Old     *Node           `json:"old"`
}

type WindowChange string

const (
	WindowChangeNew            WindowChange = "new"
	WindowChangeClose          WindowChange = "close"
	WindowChangeFocus          WindowChange = "focus"
	WindowChangeTitle          WindowChange = "title"
	WindowChangeFullscreenMode WindowChange = "fullscreen_mode"
	WindowChangeMove           WindowChange = "move"
	WindowChangeFloating       WindowChange = "floating"
	WindowChangeUrgent         WindowChange = "urgent"
	WindowChangeMark           WindowChange = "mark"
)

type WindowEvent struct {
	Change    WindowChange `json:"change"`
	Container Node         `json:"container"`
}

type ShutdownChange string

const ShutdownChangeExit ShutdownChange = "exit"

type ShutdownEvent struct {
	Change ShutdownChange `json:"change"`
}

type RunCommandReply struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
