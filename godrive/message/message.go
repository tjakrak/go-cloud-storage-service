package message

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
	// File *os.File
	Filename string
}

type Message struct {
	Head MessageHeader
	Name string //Capitalize means public
	Body string
}

/* constructor */
func New(ty MessageType, size int, fileName string) *Message { // return a pointer to a message, without pointer it's ectra copy
	//... various prep work, sending the header ...
	// open file pass it to io.Copy
	// file, err := os.OpenFile("test.txt", os.O_RDONLY, 0666)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// io.Copy(conn, file)

	head := MessageHeader{size, ty, fileName}
	//head := MessageHeader{size, ty, file}
	msg := Message{head, "GoDrive", "Hello world!"}

	return &msg
}

/* Reciever function */
func (m *Message) Print() {
	fmt.Println(m)
}

func (m *Message) Send(conn net.Conn) error {
	file, err := os.OpenFile(m.Head.Filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// prefix the send with a size
	// create the buffered writer ourselves so gob doesn't do it
	bconn := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(bconn)
	err2 := encoder.Encode(m)
	if err2 != nil {
		log.Fatalln(err2.Error())
	}
	// ... various prep work, sending the header ...
	// open file pass it to io.Copy
	sz, err := io.Copy(bconn, file)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("File size: %d", sz)
	// ensure all data is written out to the socket
	bconn.Flush()
	return err
}
