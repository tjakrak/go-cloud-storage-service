package main

//import (
//	"godrive/message"
//	"log"
//	"net"
//	"os"
//)


import (
//	"bufio"
//	"encoding/gob"
	"godrive/message"
//	"io"
	"log"
	"net"
	"os"
)

func check(e error) {
	if e != nil {
		log.Fatalln(e.Error())
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	msg := &message.Message{}
    msg.Head.Size = 31;
	if _, err := os.Stat("./storage"); err != nil {
		if os.IsNotExist(err) {
			err2 := os.Mkdir("./storage", 0755)
			check(err2)
		}
	}

	os.Chdir("./storage")
	newDir, err := os.Getwd()
	check(err)
	log.Printf("Current Working Directory: %s\n", newDir)

	// file, err := os.OpenFile(msg.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	// check(err)
	// log.Printf("Header size: %d\n", msg.Head.Size)
	// bytes, err := io.CopyN(file, bconn, msg.Head.Size)
	// check(err)
	// log.Printf("New file size: %d\n", bytes)
	// check(err)
    msg.Head.Type = 1
	if msg.Head.Type == 0 {
		msg.Get(conn)
	} else {
        log.Printf("%T\n", conn)
		msg.Put(conn)
	}

//    file, err := os.OpenFile("/home/rgtjakrakartadinata/P4-go-away/godrive/storage/test.txt", os.O_RDONLY, 0666)
//
//    if err != nil {
//        log.Fatalln(err.Error())
//	}
//
//	// prefix the send with a size
//	// create the buffered writer ourselves so gob doesn't do it
//	bconn := bufio.NewWriter(conn)
//	encoder := gob.NewEncoder(bconn)
//	err2 := encoder.Encode(msg)
//	check(err2)
//
//	// open file pass it to io.Copy
//	sz, err := io.Copy(bconn, file)
//	check(err)
//	log.Printf("File size: %d", sz)
	// ensure all data is written out to the socket
//	bconn.Flush()

}

func main() {
	listener, err := net.Listen("tcp", "192.168.122.215:9998")
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	for {
		if conn, err := listener.Accept(); err == nil {
			log.Println("handling connection...")
			go handleConnection(conn)
		}
	}
}
