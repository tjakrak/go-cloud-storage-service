package message

import (
	"encoding/gob"
	"fmt"
	"net"
)

// enumeration
type MessageType int

const (
	StorageRequest   MessageType = iota // 0 auto incrementing variable (1 + the last)
	RetrievalRequest                    // 1
	SearchRequest                       // 2
)

type MessageHeader struct {
	Size int
	Type MessageType
}

type Message struct {
	Head MessageHeader
	Name string //Capitalize means public
	Body string
}

/* constructor */
func New(ty MessageType, size int) *Message { // return a pointer to a message, without pointer it's ectra copy
	head := MessageHeader{size, ty}
	msg := Message{head, "GoDrive", "Hello world!"}
	return &msg
}

/* Reciever function */
func (m *Message) Print() {
	fmt.Println(m)
}

func (m *Message) Send(conn net.Conn) error {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(m)
	return err
}
