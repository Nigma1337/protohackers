package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
	WELCOME   = "welcome to chat thingy! What's your name?\n"
)

type Connection struct {
	Conn net.Conn
	Name string
}

var is_alphanumeric = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	connections := []Connection{}
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, &connections)
	}
}
func handleRequest(conn net.Conn, connections *[]Connection) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	conn.Write([]byte(WELCOME))
	name_bytes, err := reader.ReadBytes(byte('\n'))
	if err != nil {
		if err != io.EOF {
			fmt.Println("failed to read data, err:", err)
		}
		fmt.Println(err)
		return
	}
	name_string := strings.TrimSuffix(string(name_bytes), "\r\n")
	name_string = strings.TrimSuffix(name_string, "\n")
	name_string = strings.TrimSuffix(name_string, "\r")
	// If name isn't alphanumeric, kill
	if !is_alphanumeric.MatchString(name_string) {
		return
	}
	var sb strings.Builder
	sb.WriteString("* The room contains: ")
	for _, conn := range *connections {
		if conn.Name == name_string {
			return
		}
		// Might as well create list over current users here
		sb.WriteString(fmt.Sprintf("%s ", conn.Name))
	}
	*connections = append(*connections, Connection{Conn: conn, Name: name_string})
	go sendAll(fmt.Sprintf("* %s has entered the room\n", name_string), name_string, connections)
	conn.Write([]byte(sb.String() + "\n"))
	for true {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			for i, conn := range *connections {
				if conn.Name == name_string {
					*connections = remove(*connections, i)
					break
				}
			}
			if err != io.EOF {
				fmt.Println("failed to read data, err:", err)
				return
			}
			bye := fmt.Sprintf("* %s has left the room\n", name_string)
			go sendAll(bye, name_string, connections)
			return
		}
		message := fmt.Sprintf("[%s] %s", name_string, string(bytes))

		go sendAll(message, name_string, connections)
	}
}

func sendAll(text string, self string, connections *[]Connection) {
	fmt.Println(text)
	for _, conn := range *connections {
		if conn.Name != self {
			conn.Conn.Write([]byte(text))
		}
	}
}

func remove(s []Connection, i int) []Connection {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
