package main

import (
	"bufio"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"godrive/message"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type RequestHandler func(net.Conn, *bufio.Reader, *message.Message)

var handlers = map[message.MessageType]RequestHandler{
	message.StorageRequest:   handlePutReq,
	message.RetrievalRequest: handleGetReq,
	message.SearchRequest:    handleSearchReq,
	message.DeleteRequest:    handleDeleteReq,
}

var userInput []string

/* Handling connections from client */
func handleConnection(conn net.Conn) {
	log.Println("handling connection...")
	defer conn.Close()
	bconn := bufio.NewReader(conn)
	decoder := gob.NewDecoder(bconn)
	msg := &message.Message{}
	decoder.Decode(msg)

	log.Println("After decode")
	if msg.Head.Type == 3 {
		err := connectToBackUpServer(msg, conn)
		if err != nil {
			note := "backup server failed during get request"
			sendMessage(note, conn)
		}
	} else if msg.Head.Type == 0 {
		timeout := 1 * time.Second
		conn2, err := net.DialTimeout("tcp", userInput[2], timeout)
		if err != nil {
			log.Println("Backup server unreachable, error: ", err)
			log.Printf("%T", conn2)
			note := "backup server failed"
			sendMessage(note, conn)
		}
	}

	changeDirectory(msg)
	log.Printf("Filename: %s", msg.Head.Filename)
	log.Printf("Type: %d", msg.Head.Type)

	header := &msg.Head
	handler := handlers[header.Type]
	if handler != nil {
		handler(conn, bconn, msg)
	} else {
		log.Println("No handler for message type: ", header.Type)
	}

	if msg.Head.Type == 0 {
		err := connectToBackUpServer(msg, conn)
		if err != nil {
			os.Remove(msg.Head.Filename)
			os.Remove(getHashFileName(msg.Head.Filename))
			note := "backup server failed during put request"
			sendMessage(note, conn)
		}
	}
}

/* Connecting to the backup server */
func connectToBackUpServer(msg *message.Message, conn net.Conn) error {
	log.Printf("DIAL COUNTER: %d ---- %d", msg.Counter, msg.Head.Type)
	if msg.Counter > 0 {
		msg.Counter = msg.Counter - 1
		path, _ := os.Getwd()
		log.Printf("New DIAL working directory: %s\n", path)
		err := dialConnection(msg)
		return err
	}
	return nil
}

/* Change directory to storage */
func changeDirectory(msg *message.Message) string {
	path, err := os.Getwd()
	msg.Check(err)
	log.Printf("Current directory: %s\n", path)
	if !(strings.HasSuffix(path, userInput[3])) {
		if _, err := os.Stat(userInput[3]); err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(userInput[3], 0755)
				msg.Check(err)
			}
		}
		os.Chdir(userInput[3])
		path, err = os.Getwd()
		msg.Check(err)
		log.Printf("New current working directory: %s\n", path)
	}
	return path
}

/* Handling put request */
func handlePutReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	path := changeDirectory(msg)
	path += "/" + msg.Head.Filename
	var note string
	log.Printf("SERVER PUT -> Path: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.OpenFile(msg.Head.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
		msg.Check(err)
		log.Printf("SERVER PUT -> Header size: %d\n", msg.Head.Size)
		bytes, err := io.CopyN(file, bconn, msg.Head.Size)
		msg.Check(err)
		log.Printf("SERVER PUT -> New file size: %d\n", bytes)
		file.Seek(0, 0)
		writeHashFile(file, msg.Head.Filename)
		note = "File " + msg.Head.Filename + " is stored\n"
		defer file.Close()
	} else {
		note = "File already exists. Please delete existing file first.\n"
	}
	sendMessage(note, conn)
}

/* Handling get request */
func handleGetReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	log.Printf("SERVER GET -> Filename: %s", msg.Head.Filename)
	path := changeDirectory(msg)
	path += "/" + msg.Head.Filename
	log.Printf("SERVER GET -> Path: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		msg.Body = "File doesn't exist"
		log.Println(msg.Body)
		return
	} else {
		if checkMatchWithHash(msg.Head.Filename) {
			fileStat, err := os.Stat(msg.Head.Filename)
			msg.Check(err)
			msg.Head.Size = fileStat.Size()
			note := msg.Head.Filename + " is downloaded"
			msg.Body = note
			msg.PutRequest(conn)
		} else {
			msg.Body = "File is corrupted, repairing from backup server"
			log.Println(msg.Body)

			if msg.Counter > 0 {
				msg.Counter = msg.Counter - 1
				err := dialConnection(msg)
				if err != nil {
					note := "backup server failed"
					sendMessage(note, conn)
					return
				}
				msg.PutRequest(conn)
			}
			return
		}
	}
}

/* Handling search request */
func handleSearchReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	path := changeDirectory(msg)
	files, err := ioutil.ReadDir(path)
	msg.Check(err)

	var returnFiles string
	r, _ := regexp.Compile(".*?" + msg.Head.Filename + ".*")
	for _, f := range files {
		if r.MatchString(f.Name()) {
			fmt.Println(f.Name())
			returnFiles += f.Name() + "\n"
		}
	}
	sendMessage(returnFiles, conn)
}

/* Handling delete request */
func handleDeleteReq(conn net.Conn, bconn *bufio.Reader, msg *message.Message) {
	path := changeDirectory(msg)
	path += "/" + msg.Head.Filename
	var note string
	log.Printf("SERVER DEL -> Path: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		note = "File doesn't exist"
	} else {
		err := os.Remove(msg.Head.Filename)
		msg.Check(err)
		checksum := getHashFileName(msg.Head.Filename)
		err = os.Remove(checksum)
		msg.Check(err)
		note = msg.Head.Filename + " is deleted"
	}
	fmt.Println(note)
	os.Chdir("..")
	sendMessage(note, conn)
}

/* Sending notification message to client */
func sendMessage(note string, conn net.Conn) {
	conn.Write([]byte(note))
}

/* Creating a connection to a backup server */
func dialConnection(msg *message.Message) error {
	log.Printf("backup server is: %s", userInput[2])
	conn, err := net.Dial("tcp", userInput[2])

	if err != nil {
		log.Printf("ERROR in dial connection:%T", err)
		return err
	}

	defer conn.Close()

	msg.Print()
	log.Println("dial connection is sending message....")
	msg.Send(conn)
	return err
}

/* Get the md5sum of the file */
func hashFile(file *os.File) string {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}
	log.Printf("Hash: %x\n", hash.Sum(nil))
	return hex.EncodeToString(hash.Sum(nil))
}

/* Write to filehash.txt */
func writeHashFile(file *os.File, fileName string) {
	hash := hashFile(file)
	fileHash := getHashFileName(fileName)
	writeToFile(hash, fileHash)
}

/* Add checksum to the end of file */
func getHashFileName(fileName string) string {
	return fileName + ".checksum"
}

/* Check if the hashes matched */
func checkMatchWithHash(fileName string) bool {
	fileHash := getHashFileName(fileName)
	storedHash, err := ioutil.ReadFile(fileHash)
	if err != nil {
		log.Fatalf("failed reading from file: %s", err)
	}
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed opening from file: %s", err)
	}
	receivedHash := hashFile(file)
	defer file.Close()
	log.Printf("Stored hash : %s\n", string(storedHash))
	log.Printf("Received hash : %s\n", receivedHash)
	if strings.Compare(string(storedHash), receivedHash) == 0 {
		return true
	} else {
		log.Println("Hashes don't match")
		return false
	}
}

/* Write to file */
func writeToFile(hash string, file string) {
	openFile(file)
	hashData := []byte(hash)
	err := ioutil.WriteFile(file, hashData, 0777)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

/* Open file */
func openFile(file string) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("failed opening file: %s", err.Error())
		return
	}
	defer f.Close()
}

func main() {
	userInput = os.Args
	log.Printf("%s", userInput[1])
	listener, err := net.Listen("tcp", userInput[1])

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	for {
		if conn, err := listener.Accept(); err == nil {
			go handleConnection(conn)
		}
	}
}
