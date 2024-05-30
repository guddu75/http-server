package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var directory string

type Request struct {
	method  string
	path    string
	headers map[string]string
	body    string
}

// func parseRe(con net.Conn) Request {

// }

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
	lines := strings.Split(content, "\r\n")
	fmt.Println("lines here", len(lines), lines)

	userAgent := strings.Split(lines[2], ": ")[1]
	resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(userAgent)) + "\r\n\r\n" + userAgent
	con.Write([]byte(resp))
}

func responseFile(con net.Conn, path string) {
	filename := strings.Split(path, "/")[2]
	filepath := fmt.Sprintf("%s/%s", directory, filename)
	_, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		response404(con)
	} else {
		content, err := os.ReadFile(filepath)

		data := string(content)

		if err != nil {
			fmt.Println("Can not open file")
		} else {
			resp := "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + fmt.Sprint(len(data)) + "\r\n\r\n" + data
			con.Write([]byte(resp))
		}
	}
}

// func postFile(con net.Conn, content string) {
// 	lines := strings.Split(content, "\r\n")
// 	fileContent :=
// }

func handleRequest(con net.Conn) {

	defer con.Close()

	buf := make([]byte, 1024)

	contentLength, _ := con.Read(buf)

	content := string(buf[:contentLength])

	fmt.Printf("content here %q", content)

	path := strings.Split(content, " ")[1]

	method := strings.Split(content, " ")[0]

	if method == "POST" {
		if strings.HasPrefix(path, "/files/") {
			// postFile(con, content)
		}
	} else {
		if path == "/" {
			response200(con)
		} else if path == "/user-agent" {
			responseUserAgent(con, content)
		} else if strings.HasPrefix(path, "/echo/") {
			responseEcho(con, path)
		} else if strings.HasPrefix(path, "/files/") {
			responseFile(con, path)
		} else {
			response404(con)
		}
	}

}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	args := os.Args

	if len(args) > 2 && args[1] == "--directory" {
		directory = args[2]
	}

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {
		con, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(con)

	}

	//response

}
