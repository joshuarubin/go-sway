package sway

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/joshuarubin/lifecycle"
	"go.uber.org/multierr"
)

type client struct {
	conn net.Conn
	path string
}

// A Client provides simple communication with the sway IPC
type Client interface {
	// Runs the payload as sway commands
	RunCommand(context.Context, string) ([]RunCommandReply, error)

	// Get the list of current workspaces
	GetWorkspaces(context.Context) ([]Workspace, error)

	// Get the list of current outputs
	GetOutputs(context.Context) ([]Output, error)

	// Get the node layout tree
	GetTree(context.Context) (*Node, error)

	// Get the names of all the marks currently set
	GetMarks(context.Context) ([]string, error)

	// Get the list of configured bar IDs
	GetBarIDs(context.Context) ([]string, error)

	// Get the specified bar config
	GetBarConfig(context.Context, string) (*BarConfig, error)

	// Get the version of sway that owns the IPC socket
	GetVersion(context.Context) (*Version, error)

	// Get the list of binding mode names
	GetBindingModes(context.Context) ([]string, error)

	// Returns the config that was last loaded
	GetConfig(context.Context) (*Config, error)

	// Sends a tick event with the specified payload
	SendTick(context.Context, string) (*TickReply, error)

	// Get the list of input devices
	GetInputs(context.Context) ([]Input, error)

	// Get the list of seats
	GetSeats(context.Context) ([]Seat, error)
}

// Option can be passed to New to specify runtime configuration settings
type Option func(*client)

// WithSocketPath explicitly sets the sway socket path so it isn't read from
// $SWAYSOCK
func WithSocketPath(socketPath string) Option {
	return func(c *client) {
		c.path = socketPath
	}
}

// New returns a Client configured to connect to $SWAYSOCK
func New(ctx context.Context, opts ...Option) (_ Client, err error) {
	c := &client{}

	for _, opt := range opts {
		opt(c)
	}

	if c.path == "" {
		c.path = strings.TrimSpace(os.Getenv("SWAYSOCK"))
	}

	if c.path == "" {
		return nil, fmt.Errorf("$SWAYSOCK is empty")
	}

	c.conn, err = (&net.Dialer{}).DialContext(ctx, "unix", c.path)

	if lifecycle.Exists(ctx) {
		lifecycle.DeferErr(ctx, c.conn.Close)
	} else {
		go func() {
			<-ctx.Done()
			_ = c.conn.Close()
		}()
	}

	return c, err
}

type payloadReader struct {
	io.Reader
}

func (r payloadReader) Close() error {
	_, err := ioutil.ReadAll(r)
	return err
}

func (c *client) recvMsg(ctx context.Context) (*message, error) {
	var h header
	err := do(ctx, func() error {
		return binary.Read(c.conn, binary.LittleEndian, &h)
	})
	if err != nil {
		return nil, err
	}

	return &message{
		Type:    h.Type,
		Payload: payloadReader{io.LimitReader(c.conn, int64(h.Length))},
	}, nil
}

func (c *client) roundTrip(ctx context.Context, t messageType, payload []byte) (*message, error) {
	if c == nil {
		return nil, fmt.Errorf("not connected")
	}

	err := do(ctx, func() error {
		err := binary.Write(c.conn, binary.LittleEndian, &header{magic, uint32(len(payload)), t})
		if err != nil {
			return nil
		}

		_, err = c.conn.Write(payload)
		return err
	})
	if err != nil {
		return nil, err
	}

	return c.recvMsg(ctx)
}

func (c *client) GetTree(ctx context.Context) (*Node, error) {
	b, err := c.roundTrip(ctx, messageTypeGetTree, nil)
	if err != nil {
		return nil, err
	}

	var n Node
	return &n, b.Decode(&n)
}

func (c *client) subscribe(ctx context.Context, events ...EventType) error {
	payload, err := json.Marshal(events)
	if err != nil {
		return err
	}

	msg, err := c.roundTrip(ctx, messageTypeSubscribe, payload)
	if err != nil {
		return err
	}

	var reply struct {
		Success bool `json:"success"`
	}

	if err = msg.Decode(&reply); err != nil {
		return err
	}

	if !reply.Success {
		return fmt.Errorf("subscribe unsuccessful")
	}

	return nil
}

func (c *client) RunCommand(ctx context.Context, command string) ([]RunCommandReply, error) {
	msg, err := c.roundTrip(ctx, messageTypeRunCommand, []byte(command))
	if err != nil {
		return nil, err
	}

	var replies []RunCommandReply
	if err = msg.Decode(&replies); err != nil {
		return nil, err
	}

	for _, reply := range replies {
		if !reply.Success {
			err = multierr.Append(err, fmt.Errorf("command %q unsuccessful: %v", command, reply.Error))
		}
	}

	return replies, err
}

func (c *client) GetWorkspaces(ctx context.Context) ([]Workspace, error) {
	msg, err := c.roundTrip(ctx, messageTypeGetWorkspaces, nil)
	if err != nil {
		return nil, err
	}

	var ret []Workspace
	return ret, msg.Decode(&ret)
}

func (c *client) GetOutputs(ctx context.Context) ([]Output, error) {
	msg, err := c.roundTrip(ctx, messageTypeGetOutputs, nil)
	if err != nil {
		return nil, err
	}

	var ret []Output
	return ret, msg.Decode(&ret)
}

func (c *client) GetMarks(ctx context.Context) ([]string, error) {
	msg, err := c.roundTrip(ctx, messageTypeGetMarks, nil)
	if err != nil {
		return nil, err
	}

	var ret []string
	return ret, msg.Decode(&ret)
}

func (c *client) GetBarIDs(ctx context.Context) ([]string, error) {
	msg, err := c.roundTrip(ctx, messageTypeGetBarConfig, nil)
	if err != nil {
		return nil, err
	}

	var ret []string
	return ret, msg.Decode(&ret)
}

func (c *client) GetBarConfig(ctx context.Context, id string) (*BarConfig, error) {
	msg, err := c.roundTrip(ctx, messageTypeGetBarConfig, []byte(id))
	if err != nil {
		return nil, err
	}

	var ret BarConfig
	return &ret, msg.Decode(&ret)
}

func (c *client) GetVersion(ctx context.Context) (*Version, error) {
	msg, err := c.roundTrip(ctx, messageTypeGetVersion, nil)
	if err != nil {
		return nil, err
	}

	var ret Version
	return &ret, msg.Decode(&ret)
}

func (c *client) GetBindingModes(ctx context.Context) ([]string, error) {
	msg, err := c.roundTrip(ctx, messageTypeGetBindingModes, nil)
	if err != nil {
		return nil, err
	}

	var ret []string
	return ret, msg.Decode(&ret)
}

func (c *client) GetConfig(ctx context.Context) (*Config, error) {
	msg, err := c.roundTrip(ctx, messageTypeGetConfig, nil)
	if err != nil {
		return nil, err
	}

	var ret Config
	return &ret, msg.Decode(&ret)
}

func (c *client) SendTick(ctx context.Context, payload string) (*TickReply, error) {
	msg, err := c.roundTrip(ctx, messageTypeSendTick, []byte(payload))
	if err != nil {
		return nil, err
	}

	var ret TickReply
	return &ret, msg.Decode(&ret)
}

func (c *client) GetInputs(ctx context.Context) ([]Input, error) {
	msg, err := c.roundTrip(ctx, messageTypeGetInputs, nil)
	if err != nil {
		return nil, err
	}

	var ret []Input
	return ret, msg.Decode(&ret)
}

func (c *client) GetSeats(ctx context.Context) ([]Seat, error) {
	msg, err := c.roundTrip(ctx, messageTypeGetSeats, nil)
	if err != nil {
		return nil, err
	}

	var ret []Seat
	return ret, msg.Decode(&ret)
}
