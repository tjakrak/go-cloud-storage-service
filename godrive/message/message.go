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

type MessageType int

const (
	StorageRequest   MessageType = iota // 0
	RetrievalRequest                    // 1
	SearchRequest                       // 2
	DeleteRequest                       // 3
)

type MessageHeader struct {
	Size     int64
	Type     MessageType
	Filename string
}

type Message struct {
	Head MessageHeader
	Name string
	Body string
}

/* constructor */
func New(ty MessageType, size int64, fileName string) *Message { // return a pointer to a message, without pointer it's extra copy
	head := MessageHeader{size, ty, fileName}
	var request string
	if head.Type == 0 {
		request = "put"
	} else if head.Type == 1 {
		request = "get"
	}
	msg := Message{head, head.Filename, request}

	return &msg
}

/* Reciever function */
func (m *Message) Print() {
	fmt.Println(m)
}

func (m *Message) Send(conn net.Conn) error {
	log.Printf("In message send. File type: %d", m.Head.Type)
	var err error
	if m.Head.Type == 0 {
		err = m.Put(conn)
		check(err)
	} else if m.Head.Type == 1 {
		log.Printf("In message send if get. File name: %s", m.Head.Filename)
		err = m.GetRequest(conn)
		check(err)
	}
	return err
}

func (m *Message) Put(conn net.Conn) error {
	file, err := os.OpenFile(m.Head.Filename, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalln(err.Error())
	}

	bconn := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(bconn)
	err2 := encoder.Encode(m)
	if err2 != nil {
		log.Fatalln(err2.Error())
	}

	sz, err := io.Copy(bconn, file)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("File size: %d", sz)

	bconn.Flush()
	return err
}

func (m *Message) Get(bconn *bufio.Reader) error {
	file, err := os.OpenFile(m.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	check(err)
	log.Printf("Header size: %d\n", m.Head.Size)
	bytes, err := io.CopyN(file, bconn, m.Head.Size)
	check(err)

	log.Printf("New file size: %d\n", bytes)
	return err
}

func (m *Message) GetRequest(conn net.Conn) error {
	bconn := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(bconn)
	err2 := encoder.Encode(m)
	if err2 != nil {
		log.Fatalln(err2.Error())
	}

	bconn.Flush()

	cconn := bufio.NewReader(conn)
	decoder := gob.NewDecoder(cconn)
	err3 := decoder.Decode(m)
	check(err3)

	file, err := os.OpenFile(m.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	check(err)
	log.Printf("Header size: %d\n", m.Head.Size)
	bytes, err := io.CopyN(file, cconn, m.Head.Size)
	check(err)

	log.Printf("New file size: %d\n", bytes)
	return err2
}

func check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}
