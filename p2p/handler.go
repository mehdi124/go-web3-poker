package p2p

import (
	"io"
)

type Handler interface {
	HandleMessage(*Message) error
}

type DefaultHandler struct{}

func (h *DefaultHandler) HandleMessage(msg *Message) error {

	b, err := io.ReadAll(msg.Payload)
	if err != nil {
		return err
	}
	fmr.Printf("handling the msg from %s:%s", msg.From, string(b))

	return nil
}
