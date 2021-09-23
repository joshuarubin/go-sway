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
	messageTypeGetInputs = 100
	messageTypeGetSeats  = 101
)

const (
	eventTypeWorkspace       messageType = 0x80000000
	eventTypeMode            messageType = 0x80000002
	eventTypeWindow          messageType = 0x80000003
	eventTypeBarConfigUpdate messageType = 0x80000004
	eventTypeBinding         messageType = 0x80000005
	eventTypeShutdown        messageType = 0x80000006
	eventTypeTick            messageType = 0x80000007
	eventTypeBarStateUpdate  messageType = 0x80000014
	eventTypeInput           messageType = 0x80000015
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
