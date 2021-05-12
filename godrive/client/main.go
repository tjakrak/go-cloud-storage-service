package main

import (
	"fmt"
	"godrive/message"
	"log"
	"net"
	"os"
)

func main() {
	// commandline processing
	// path manipulation

	userInput := os.Args
	fmt.Printf("%s\n", userInput)
	fmt.Printf("%s\n", userInput[1])
	fmt.Printf("%d\n", len(userInput))

	conn, err := net.Dial("tcp", userInput[1])
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	defer conn.Close()
	var msg *message.Message

	// figure out put or get userInput[2]
	if userInput[2] == "put" {
		fileStat, err := os.Stat(userInput[3])
		if err != nil {
			log.Println("Hey error in the os.Stat CLIENT")
			log.Fatalln(err.Error())
			return
		}
		fileSize := fileStat.Size()
		log.Printf("File Size: %d\n", fileSize)
		msg = message.New(0, fileSize, userInput[3]) //use os.stat
	} else if userInput[2] == "get" {
		msg = message.New(1, 17, userInput[3])
		// } else if userInput[2] == "search" {
		// 	msg = message.New(2, fileSize, userInput[3])
		// } else if userInput[2] == "delete" {
		// 	msg = message.New(3, fileSize, userInput[3])
	} else {
		log.Fatalln(err.Error())
	}

	msg.Print()
	msg.Send(conn) // pass in our connection
}
