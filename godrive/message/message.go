package message

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
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
	//... various prep work, sending the header ...
	// open file pass it to io.Copy

	return &msg
}

/* Reciever function */
func (m *Message) Print() {
	fmt.Println(m)
}

func (m *Message) Send(conn net.Conn) error {
	// prefix the send with a size
	// create the buffered writer ourselves so gob doesn't do it
	bconn := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(bconn)
	err := encoder.Encode(m)
	//... various prep work, sending the header ...
	// open file pass it to io.Copy
	sz, err := io.Copy(bconn, file)
	// ensure all data is written out to the socket
	bconn.Flush()
	return err
}
