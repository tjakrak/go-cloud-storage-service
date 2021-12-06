# Cloud Storage Service
Author: Marisa Tania, Ryan Tjakrakartadinata\
Professor: Matthew Malensek\
See project spec here: https://www.cs.usfca.edu/~mmalensek/cs521/assignments/project-4.html

## About GoDrive
GoDrive is a cloud storage system similar to Dropbox or Google Drive, with resilient, replicated backend servers and a command line client application. This project uses the go standard library and socket programming.

### GoDrive Design
![GoDrive](https://user-images.githubusercontent.com/60201466/119010812-49fcfd80-b949-11eb-87fc-2c9f3837cd1f.jpg)


### How to use GoDrive

To set up the server: 
```bash
go run ./server/main.go host:port host:port storageA
                        SERVER 1   SERVER 2 STORAGE_NAME
```

To set up the client: 
Client can pick to input the port for server 1 or server 2
```bash
go run ./client/main.go host:port fileName
                        SERVER 1  FILE NAME
```

To compile and run:

Main Server
```bash
go run ./server/main.go localhost:7777 localhost:7778 storageA
go run ./server/main.go 192.168.122.215:7777 192.168.122.215:7778 storageA
```

Backup Server
```bash
go run ./server/main.go localhost:7778 localhost:7777 storageB
go run ./server/main.go 192.168.122.215:7778 192.168.122.215:7777 storageB
```

Client
```bash
go run client/main.go 192.168.122.212:7777 put ./file.txt
go run client/main.go 192.168.122.212:7777 get file.txt
go run client/main.go 192.168.122.212:7777 search file.txt
go run client/main.go 192.168.122.212:7777 delete file.txt
```

### GoDrive Components
GoDrive has two components:
- <b>Storage Server</b>: handles storage/retrieval/search operations. Will replicate files to another storage server instance.
- <b>Client</b>: can send requests to any of the storage servers.

### GoDrive Features
Specific features of the system include:
- <b>Storage</b>: Store any type of file (text, images, binaries, and so on). Given enough disk space (and time to transfer the data), GoDrive supports arbitrarily large file sizes.
- <b>Retrieval</b>: beyond usual file retrievals with get, GoDrive is able to search and list the files in the system
- <b>Scalability</b>: concurrent storage and retrieval operations is supported, as well as handling multiple clients.
- <b>Replication</b>: backend servers ensure that all files are replicated for fault tolerance. If a backend server goes down, GoDrive is able to contact a replica for the file.
- <b>Resiliency</b>: GoDrive is resilient to disk or memory failures. If a file is stored on a disk and gets corrupted, GoDrive can detect the corruption, retrieve a replica, and repair the file.

### Storage/Retrieval Operations
Both the server and client in GoDrive support a variety of messages that influence behavior. Here are the operations:
```bash
Options:
    * put fileName      Store file
    * get fileName      Retrieve file
    * delete fileName   Delete file
    * search string     (note that this string could be blank to search for all files)
```
To ensure the system is trustworthy, GoDrive acknowledge each of these operations as either successful or a failure. 

### Included Files and Directories
There are several files included. These are:
   - <b>Makefile</b>: Including to compile and run the program.
   - <b>server/main.go</b>: The server driver that listens and responses to client requests
   - <b>client/main.go</b>: The client driver that send out message requests
   - <b>message.go</b>: Shared put, get, delete and search message requests
 
 There are also directories to store shared files:
   - <b>storageA</b>: for the main server
   - <b>storageB</b>: for the backend server

### Program Output

Server
```bash
[mtania@nemo godrive]$ go run ./server/main.go 192.168.122.212:7778 192.168.122.212:7777 storageA
2021/05/14 00:23:18 192.168.122.212:7778
2021/05/14 00:24:43 handling connection...
2021/05/14 00:24:43 DIAL COUNTER: 1 ---- 1
2021/05/14 00:24:43 File name: file.txt
2021/05/14 00:24:43 Type: 1
2021/05/14 00:24:43 File size: 29kB 
```

Client
```bash
2021/05/15 00:24:43 Current directory: /home/mtania/P4-go-away/godrive/storageA
2021/05/15 00:24:43 SERVER GET -> Path: /home/mtania/P4-go-away/godrive/storageA/file.txt
2021/05/15 00:24:43 File is downloaded
```
