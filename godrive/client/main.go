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
	// conn, err := net.Dial("tcp", ":9999") // connect to localhost port 9999
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	// figure out put or get userInput[2]

	// something := message.SearchRequest
	// m := message.Message{Name: "GoDrive"}
	// fmt.Println(m, something)
	defer conn.Close()
	msg := message.New(message.SearchRequest, 300, userInput[3])
	msg.Print()
	msg.Send(conn) // pass in our connection

	// move to the constructor open the file
	// file, err := os.OpenFile("test.txt", os.O_RDONLY, 0666)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// 	return
	// }
	// io.Copy(conn, file)
}
