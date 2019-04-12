package sway

import (
	"context"
	"encoding/json"
	"io"
)

type header struct {
	Magic  [6]byte
	Length uint32
	Type   messageType
}

type message struct {
	Type    messageType
	Payload io.ReadCloser
}

func (m message) Decode(v interface{}) error {
	defer m.Payload.Close()
	return json.NewDecoder(m.Payload).Decode(v)
}

type messageType uint32

const (
	messageTypeRunCommand messageType = iota
	messageTypeGetWorkspaces
	messageTypeSubscribe
	messageTypeGetOutputs
	messageTypeGetTree
	messageTypeGetMarks
	messageTypeGetBarConfig
	messageTypeGetVersion
	messageTypeGetBindingModes
	messageTypeGetConfig
	messageTypeSendTick
	messageTypeSync
)

const (
	eventReplyTypeWorkspace       messageType = 0x80000000
	eventReplyTypeMode            messageType = 0x80000002
	eventReplyTypeWindow          messageType = 0x80000003
	eventReplyTypeBarConfigUpdate messageType = 0x80000004
	eventReplyTypeBinding         messageType = 0x80000005
	eventReplyTypeShutdown        messageType = 0x80000006
	eventReplyTypeTick            messageType = 0x80000007
	eventReplyTypeBarStatusUpdate messageType = 0x80000014
)

var magic = [6]byte{'i', '3', '-', 'i', 'p', 'c'}

func do(ctx context.Context, fn func() error) error {
	done := make(chan struct{})
	var err error
	go func() {
		err = fn()
		close(done)
	}()

	select {
	case <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
