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

// type MessageHandler func(*message.Header, *message.Message)
// var handlers = map[message.MessageType]MessageHandler{
// message.PutRequest: handlePutReq,
// message.GetRequest: handleGetReq,
// }

const (
	StorageRequest   MessageType = iota // 0 auto incrementing variable (1 + the last)
	RetrievalRequest                    // 1
	SearchRequest                       // 2
	DeleteRequest                       //--> use in the server to check what kind of connection the server get
)

type MessageHeader struct {
	Size     int64
	Type     MessageType
	Filename string
}

type Message struct {
	Head MessageHeader // definitely keep the header
	Name string        //Capitalize means public
	Body string
}

/* constructor */
func New(ty MessageType, size int64, fileName string) *Message { // return a pointer to a message, without pointer it's extra copy
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
	var err error
    fmt.Printf("MSGTYPE: %d\n", m.Head.Type)
	if m.Head.Type == 0 {
		err = m.Put(conn)
	} else if m.Head.Type == 1 {
		err = m.Get(conn)
	}
	return err
}

func (m *Message) Put(conn net.Conn) error {

    file, err := os.OpenFile("test.txt", os.O_RDONLY, 0666)
//	file, err := os.OpenFile(m.Head.Filename, os.O_RDONLY, 0666)
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

func (m *Message) Get(conn net.Conn) error {
	bconn := bufio.NewReader(conn)
	decoder := gob.NewDecoder(bconn)
	err := decoder.Decode(m)
	check(err)

	file, err := os.OpenFile(m.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	check(err)
	log.Printf("Header size: %d\n", m.Head.Size)
	bytes, err := io.CopyN(file, bconn, m.Head.Size)
	check(err)
	log.Printf("New file size: %d\n", bytes)
	return err
}

func check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}
