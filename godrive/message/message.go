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
		err = m.Get(conn)
		check(err)
	}
	return err
}

func (m *Message) Put(conn net.Conn) error {
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
	// bconn := bufio.NewWriter(conn)
	// encoder := gob.NewEncoder(bconn)
	// log.Printf("MSG Before encoder: Message header filename: %s", m.Head.Filename)
	// err := encoder.Encode(m)
	// log.Printf("MSG After encoder: Message header filename: %s", m.Head.Filename)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }

	log.Printf("In message get function. File name: %s", m.Head.Filename)
	bconn2 := bufio.NewReader(conn)
	decoder := gob.NewDecoder(bconn2)
	log.Printf("Before decoder In message get function. File name: %s", m.Head.Filename)
	err := decoder.Decode(m)
	// check(err)
	// fmt.Println(m)
	if err != nil {
		log.Panic("ERROR ALERT")
		log.Fatalln(err.Error())
	}
	log.Printf("After decoder In message get function. File name: %s", m.Head.Filename)

	file, err := os.OpenFile(m.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	check(err)
	log.Printf("Header size: %d\n", m.Head.Size)
	// bytes, err := io.Copy(bconn, file)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	bytes, err := io.CopyN(file, bconn2, m.Head.Size)
	check(err)
	log.Printf("New file size: %d\n", bytes)
	return err
}

func check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}
