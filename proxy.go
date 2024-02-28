// proxy.go
// sdlewis
// 012134058
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Client is struct to hold socket , data channel and the IP of the client
type Client struct {
	socket   net.Conn
	data     chan []byte
	clientIP string
}

var goroutineCounter int

// UriParts holds the uri,hostname,pathname and the port
type UriParts struct {
	uri      string
	hostname string
	pathname string
	port     int
}

// Mutex is used to avoid race conditions
var mut sync.Mutex

func filterNewLines(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case 0x000A, 0x000B, 0x000C, 0x000D, 0x0085, 0x2028, 0x2029:
			return -1
		default:
			return r
		}
	}, s)
}

// Function to get the ID of the current goroutine
func getGoroutineID() int {
	mut.Lock()
	goroutineCounter++
	id := goroutineCounter
	mut.Unlock()
	return id
}

// this function is used to extract parts of the Uri and validate the Uri
func (uri *UriParts) parseURI() bool {
	if !strings.HasPrefix(uri.uri, "http://") && !strings.HasPrefix(uri.uri, "https://") {
		fmt.Println(uri.uri)
		uri.hostname = ""
		return false
	}
	var testsplit = strings.Split(uri.uri, "/")
	if strings.Contains(testsplit[2], ":") {
		var hostsplit = strings.Split(testsplit[2], ":")
		testsplit[2] = hostsplit[0]
		hostsplit[1] = string(bytes.Trim([]byte(filterNewLines(hostsplit[1])), "\x00"))
		port, err := strconv.Atoi(hostsplit[1])
		if err != nil {
			fmt.Println(err)
			return false
		}
		uri.port = port
	} else {
		if strings.HasPrefix(uri.uri, "https://") {
			uri.port = 443
		} else {
			uri.port = 80
		}
	}
	uri.hostname = testsplit[2]
	var paths = strings.SplitAfterN(uri.uri, "/", 4)
	if len(paths) >= 4 {
		uri.pathname = string(bytes.Trim([]byte(filterNewLines(paths[3])), "\x00"))
	} else {
		uri.pathname = ""
	}
	return true
}

// logic to execute a REST request
func executeHTTPRequest(uri *UriParts, resultChannel chan []byte) {
	url := fmt.Sprintf("http://%s:%d/%s", uri.hostname, uri.port, uri.pathname)

	//if we get an error while making a http request, following block of code will execute
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in executing the HTTP request:", err)
		resultChannel <- []byte(fmt.Sprintf("Error while  executing  the HTTP request: %s", err))
		return
	}
	defer resp.Body.Close()
	//if we get an error while reading the responses following block of code will execute
	bodyBuffer := new(bytes.Buffer)
	_, err = bodyBuffer.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Error reading HTTP response body:", err)
		resultChannel <- []byte(fmt.Sprintf("Error reading HTTP response body: %s", err))
		return
	}
	resultChannel <- bodyBuffer.Bytes() //if we get a proper response , we flush it into the channel
}

// function to write to a file
func writeLogEntryToFile(logEntry, filename string, goroutineID int) {
	mut.Lock() //locking the resource aka file
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print(err)
		mut.Unlock()
		return
	}
	defer file.Close()
	logEntry += fmt.Sprintf(" GoroutineID: %d\n", goroutineID)
	logEntry += "\n"
	logEntry = strings.ReplaceAll(logEntry, "\n", " ")
	if _, err := file.WriteString(logEntry + "\n"); err != nil {
		fmt.Print(err)
	}
	mut.Unlock() //unclock after doing necessary operation
}

func (client *Client) receive() {
	defer client.socket.Close()
	reader := bufio.NewReader(client.socket) //create a buffer reader that reads from the network connection
	client.socket.Write([]byte("\n"))
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err)
			break
		}
		message = strings.TrimRight(message, "\r\n")
		uri := &UriParts{uri: message}
		if uri.parseURI() {

			resultChannel := make(chan []byte) // creating a channel so that the go routines can communicate with each other through it
			//go executeCurlCommand(curlCommand, resultChannel)
			go executeHTTPRequest(uri, resultChannel)
			result := <-resultChannel //listen to the channel
			//formatting the log
			currentDate := time.Now().Format("Mon 02 Jan 2006 15:04:05 PST")
			size := len(result)
			logEntry := fmt.Sprintf("Date: %s Client IP: %s URL: %s Size: %d ", currentDate, client.clientIP, message, size)
			goroutineID := getGoroutineID()
			go writeLogEntryToFile(logEntry, "proxy.log", goroutineID)

			// Echo the received message back to the client
			result = append(result, '\n')
			result = bytes.ReplaceAll(result, []byte("\n"), []byte("\r\n"))
			client.socket.Write(result)
		} else {
			fmt.Println("Failed to parse URI.")
		}
	}
	fmt.Println("Socket Closed")
}

// function to create  a server and listen on a specific  port
func startServer(port string) {
	listener, error := net.Listen("tcp", ":"+port)
	if error != nil {
		fmt.Println(error)
		return
	}
	for {
		connection, error := listener.Accept()
		if error != nil {
			fmt.Println(error)
		} else {
			clientIP := connection.RemoteAddr().(*net.TCPAddr).IP.String()
			client := &Client{socket: connection, data: make(chan []byte), clientIP: clientIP}
			go client.receive() //create  a network connection
		}

	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) >= 1 {
		startServer(argsWithoutProg[0])
	} else {
		startServer("1234")
	}
}
