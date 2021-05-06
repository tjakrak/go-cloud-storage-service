# Project 4: Cloud Storage Service
Author: Marisa Tania, Ryan Tjakrakartadinata
Professor: Matthew Malensek
See project spec here: https://www.cs.usfca.edu/~mmalensek/cs521/assignments/project-4.html

## About CSS
CSS is a cloud storage system similar to Dropbox or Google Drive, with resilient, replicated backend servers and a command line client application. This project uses the go standard library and socket programming.

### CSS Design
TBD discuss your design and the logic behind the decisions you made

### How to use CSS
TBD detailed instructions on how to set up and use CSS

### CSS Features
Specific features of the system include:
- <b>Storage</b>: Store any type of file (text, images, binaries, and so on). Given enough disk space (and time to transfer the data), CSS supports arbitrarily large file sizes.
- <b>Retrieval</b>: beyond usual file retrievals with get, CSS is able to search and list the files in the system
- <b>Scalability</b>: concurrent storage and retrieval operations is supported, as well as handling multiple clients.
- <b>Replication</b>: backend servers ensure that all files are replicated for fault tolerance. If a backend server goes down, CSS is able to contact a replica for the file.
- <b>Resiliency</b>: CSS is resilient to disk or memory failures. If a file is stored on a disk and gets corrupted, CSS can detect the corruption, retrieve a replica, and repair the file.

### CSS Components
CSS has two components:
- <b>Storage Server</b>: handles storage/retrieval/search operations. Will replicate files to another storage server instance.
- <b>Client</b>: can send requests to any of the storage servers.

### Storage/Retrieval Operations
Both the server and client in CSS support a variety of messages that influence behavior. Here are the operations:
```bash
Options:
    * put fileName      Store file
    * get fileName      Retrieve file
    * delete fileName   Delete file
    * search string     (note that this string could be blank to search for all files)
```
The task list, hardware information, system information, and task information can all be turned on/off with the command line options. By default, all of them are displayed.

### Included Files
There are several files included. These are:
   - <b>Makefile</b>: Including to compile and run the program.
   - <b>css.c</b>: The program driver

There are also header files for each of these files.


To compile and run:

```bash
make
./css
```


### Program Output
```bash
$ ./css
inspector.c:133:main(): Options selected: hardware system task_list task_summary

System Information
------------------
Hostname: lunafreya
Kernel Version: 4.20.3-arch1-1-ARCH
Uptime: 1 days, 11 hours, 49 minutes, 56 seconds

Hardware Information
------------------
CPU Model: AMD EPYC Processor (with IBPB)
Processing Units: 2
Load Average (1/5/15 min): 0.00 0.00 0.00
CPU Usage:    [--------------------] 0.0%
Memory Usage: [#-------------------] 9.5% (0.1 GB / 1.0 GB)

Task Information
----------------
Tasks running: 88
Since boot:
    Interrupts: 2153905
    Context Switches: 3678668
    Forks: 38849

  PID |        State |                 Task Name |            User | Tasks 
------+--------------+---------------------------+-----------------+-------
    1 |     sleeping |                   systemd |               0 | 1 
    2 |     sleeping |                  kthreadd |               0 | 1 
    3 |         idle |                    rcu_gp |               0 | 1 
    4 |         idle |                rcu_par_gp |               0 | 1 
    6 |         idle |      kworker/0:0H-kblockd |               0 | 1 

```
