package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	OK        = "200 OK"
	CREATED   = "201 Created"
	NOT_FOUND = "404 Not Found"
)

var directory string

type Request struct {
	method      string
	path        string
	httpVersion string
	host        string
	headers     map[string]string
	body        string
}

type Response struct {
	httpVersion string
	resCode     string
	headers     map[string]string
	body        string
}

func parseRequest(con net.Conn) *Request {
	buf := make([]byte, 1024)

	contentLength, _ := con.Read(buf)

	content := string(buf[:contentLength])

	req := new(Request)

	fmt.Println("content here ", content)

	lines := strings.Split(content, "\r\n")

	// log.Print("lines", lines, len(lines))

	for i, line := range lines {
		fmt.Println(i, line)
	}

	firstLine := strings.Split(lines[0], " ")

	log.Print("firstLine", firstLine)

	req.method = firstLine[0]

	req.path = firstLine[1]

	req.httpVersion = firstLine[2]

	if len(lines) >= 2 {
		req.host = strings.Split(lines[1], ": ")[1]
	}

	if len(lines) >= 3 {
		headers := strings.Split(lines[2], ": ")

		log.Print("headers", headers, len(headers))

		req.headers = make(map[string]string)
		if len(headers) > 1 {
			for i := 0; i < len(headers); i += 2 {
				req.headers[headers[i]] = headers[i+1]
			}
		}

		log.Print("here")
	}

	//body present or not

	if len(lines) >= 5 {
		req.body = lines[4]
	}

	// fmt.Println("request : ", req)

	return req
}

func encode(msg string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(msg)); err != nil {
		log.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("compressed msg", b.String())
	return b.String()
}

func responseEcho(con net.Conn, req Request) {
	msg := strings.Split(req.path, "/")[2]
	enc := ""
	// if req.headers["Accept-Encoding"] == "gzip" {
	// 	enc = "Content-Encoding: gzip\r\n"
	// }
	encodingMethods := strings.Split(req.headers["Accept-Encoding"], ", ")

	for _, method := range encodingMethods {
		if method == "gzip" {
			enc = "Content-Encoding: gzip\r\n"
			msg = encode(msg)
			break
		}
	}

	resp := "HTTP/1.1 200 OK\r\n" + enc + "Content-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(msg)) + "\r\n\r\n" + msg
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

func responseUserAgent(con net.Conn, req Request) {
	// lines := strings.Split(content, "\r\n")
	// fmt.Println("lines here", len(lines), lines)

	userAgent := req.headers["User-Agent"]
	resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + fmt.Sprint(len(userAgent)) + "\r\n\r\n" + userAgent
	con.Write([]byte(resp))
}

func responseFile(con net.Conn, req Request) {
	filename := strings.Split(req.path, "/")[2]
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

func postFile(con net.Conn, req Request) {
	filename := strings.Split(req.path, "/")[2]
	filepath := fmt.Sprintf("%s/%s", directory, filename)

	err := os.WriteFile(filepath, []byte(req.body), 0666)

	if err != nil {
		log.Fatal(err)
	}

	con.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
}

func handleRequest(con net.Conn) {

	defer con.Close()

	req := parseRequest(con)

	log.Print("Request parsed successfully")

	if req.method == "POST" {
		if strings.HasPrefix(req.path, "/files/") {
			postFile(con, *req)
		}
	} else if req.method == "GET" {
		if req.path == "/" {
			response200(con)
		} else if req.path == "/user-agent" {
			responseUserAgent(con, *req)
		} else if strings.HasPrefix(req.path, "/echo/") {
			responseEcho(con, *req)
		} else if strings.HasPrefix(req.path, "/files/") {
			responseFile(con, *req)
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
