package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func responseEcho(con net.Conn, path string) {
	msg := strings.Split(path, "/")[2]
	resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(msg)) + "\r\n\r\n" + msg
	con.Write([]byte(resp))
}

func response200(con net.Conn) {
	con.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}

/*
HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nfoobar/1.2.3
*/

func response404(con net.Conn) {
	con.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

func responseUserAgent(con net.Conn, content string) {
	lines := strings.Split(content, "/r/n")
	fmt.Println("lines here", lines)
	userAgent := strings.Split(lines[2], ": ")[1]
	resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(userAgent)) + "\r\n\r\n" + userAgent
	con.Write([]byte(resp))
}

func handleRequest(con net.Conn) {

	defer con.Close()
	buf := make([]byte, 1024)

	contentLength, _ := con.Read(buf)

	content := string(buf[:contentLength])

	fmt.Println("content here", content)

	path := strings.Split(content, " ")[1]

	if path == "/" {
		response200(con)
	} else if path == "/user-agent" {
		responseUserAgent(con, content)
	} else if strings.HasPrefix(path, "/echo/") {
		responseEcho(con, path)
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
