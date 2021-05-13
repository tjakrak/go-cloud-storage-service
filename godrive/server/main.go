package main

//import (
//	"godrive/message"
//	"log"
//	"net"
//	"os"
//)


import (
	"bufio"
	"encoding/gob"
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
    bconn := bufio.NewReader(conn)
    decoder := gob.NewDecoder(bconn)
	msg := &message.Message{}
    decoder.Decode(msg)

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


    log.Printf("Filename: %s\n", msg.Head.Filename)
    log.Printf("Type: %d\n", msg.Head.Type)

    if msg.Head.Type == 0 {
        msg.Get(bconn)
	} else if msg.Head.Type == 1 {
	    fileStat, err := os.Stat(msg.Head.Filename)
	    if err != nil {
		    log.Fatalln(err.Error())
		    return
        }
        msg.Head.Size = fileStat.Size()
		msg.Put(conn)
	}

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
