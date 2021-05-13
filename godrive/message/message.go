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

/* Constructor */
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

/* Sending connection based on request type */
func (m *Message) Send(conn net.Conn) error {
	var err error
	if m.Head.Type == 0 {
		err = m.Put(conn)
		check(err)
	} else if m.Head.Type == 1 {
		err = m.GetRequest(conn)
		check(err)
	}
	return err
}

/* PutRequest storing file */
func (m *Message) Put(conn net.Conn) error {

	file, err := os.OpenFile("test.txt", os.O_RDONLY, 0666)
	//	file, err := os.OpenFile(m.Head.Filename, os.O_RDONLY, 0666)
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
	log.Printf("Storing File size: %d", sz)

	bconn.Flush()
	return err
}

/* GetRequest to retrieve file */
func (m *Message) GetRequest(conn net.Conn) error {
	bconn := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(bconn)
	err := encoder.Encode(m)
	check(err)

	bconn.Flush()

	cconn := bufio.NewReader(conn)
	decoder := gob.NewDecoder(cconn)
	err = decoder.Decode(m)
	check(err)

	file, err := os.OpenFile(m.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	check(err)
	log.Printf("MSG GetRequest -> Header size: %d\n", m.Head.Size)
	bytes, err := io.CopyN(file, cconn, m.Head.Size)
	check(err)

	log.Printf("MSG GetRequest -> New file size: %d\n", bytes)
	return err
}

/* Check error */
func check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}
