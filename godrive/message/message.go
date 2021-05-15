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
	StorageRequest MessageType = iota
	RetrievalRequest
	SearchRequest
	DeleteRequest
)

type DialCounter int

type MessageHeader struct {
	Size     int64
	Type     MessageType
	Filename string
}

type Message struct {
	Head MessageHeader
	Counter DialCounter
	Body string
}

/* Constructor */
func New(ty MessageType, size int64, fileName string) *Message { // return a pointer to a message, without pointer it's extra copy
	head := MessageHeader{size, ty, fileName}
	msg := Message{head, 1, ""}
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
		err = m.PutRequest(conn)
	} else if m.Head.Type == 1 {
		err = m.GetRequest(conn)
	} else if m.Head.Type == 2 {
		err = m.SearchRequest(conn)
	} else if m.Head.Type == 3 {
		err = m.DeleteRequest(conn)
	}
	m.Check(err)
	return err
}

/* Deleting file */
func (m *Message) DeleteRequest(conn net.Conn) error {
	err := m.setEncoder(conn)
	m.Check(err)
	return err
}

/* Deleting file */
func (m *Message) SearchRequest(conn net.Conn) error {
	err := m.setEncoder(conn)
	m.Check(err)
	return err
}

/* Encoding */
func (m *Message) setEncoder(conn net.Conn) error {
	bconn := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(bconn)
	err := encoder.Encode(m)
	m.Check(err)
	bconn.Flush()
	return err
}

/* PutRequest storing file */
func (m *Message) PutRequest(conn net.Conn) error {
	file, err := os.OpenFile(m.Head.Filename, os.O_RDONLY, 0666)
	m.Check(err)

	bconn := bufio.NewWriter(conn)
	encoder := gob.NewEncoder(bconn)
	err = encoder.Encode(m)
	m.Check(err)

	sz, err := io.Copy(bconn, file)
	m.Check(err)
	log.Printf("Storing file size: %d", sz)

	bconn.Flush()
	return err
}

/* GetRequest to retrieve file */
func (m *Message) GetRequest(conn net.Conn) error {
	m.setEncoder(conn)
	cconn := bufio.NewReader(conn)
	decoder := gob.NewDecoder(cconn)
	err := decoder.Decode(m)
	if err != nil {
		fmt.Println("File doesn't exist")
		return err
	}

	file, err := os.OpenFile(m.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	m.Check(err)

	log.Printf("MSG GetRequest -> Header size: %d\n", m.Head.Size)
	bytes, err := io.CopyN(file, cconn, m.Head.Size)
	m.Check(err)

	log.Printf("MSG GetRequest -> New file size: %d\n", bytes)
	fmt.Println(m.Body)
	return err
}

/* Check error */
func (m *Message) Check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}
