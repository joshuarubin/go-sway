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

	"go.uber.org/multierr"
)

type client struct {
	conn net.Conn
	path string
}

type Client interface {
	GetTree(context.Context) (*Node, error)
	RunCommand(context.Context, string) ([]RunCommandReply, error)
	Close() error
}

func New(ctx context.Context) (_ Client, err error) {
	c := &client{}

	if c.path = strings.TrimSpace(os.Getenv("SWAYSOCK")); c.path == "" {
		return nil, fmt.Errorf("$SWAYSOCK is empty")
	}

	c.conn, err = (&net.Dialer{}).DialContext(ctx, "unix", c.path)
	return c, err
}

func (c *client) Close() error {
	return c.conn.Close()
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
