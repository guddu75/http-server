package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func response200(con net.Conn) {
	con.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}

func response404(con net.Conn) {
	con.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

func handleRequest(con net.Conn) {

	defer con.Close()
	buf := make([]byte, 1024)

	contentLength, _ := con.Read(buf)

	content := string(buf[:contentLength])

	fmt.Println(content)

	path := strings.Split(content, " ")[1]

	if path == "/" {
		response200(con)
	} else {
		response404(con)
	}

}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	con, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handleRequest(con)
	//response

}
