package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"os"
)

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type Result struct {
	Method  string `json:"method"`
	IsPrime bool   `json:"prime"`
}

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
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}
func handleRequest(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for true {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to read data, err:", err)
			}
			fmt.Println(err)
			return
		}

		var dat map[string]interface{}
		err = json.Unmarshal(bytes, &dat)
		if err != nil {
			fmt.Println("what")
			return
		}
		if dat["method"] == nil {
			return
		}
		var method string
		switch dat["method"].(type) {
		case string:
			method = dat["method"].(string)
		default:
			return
		}
		if method != "isPrime" {
			fmt.Println("wrong method")
			return
		}
		fmt.Println(dat)
		if dat["number"] == nil {
			return
		}
		var num float64
		switch dat["number"].(type) {
		case int:
			num = dat["number"].(float64)
		case float64:
			num = dat["number"].(float64)
		default:
			return
		}
		// Floating numbers cannot be prime, as prime numbers are defined as integers
		if !isIntegral(num) {
			fmt.Println("Is not intergral!")
			return
		}
		numInt := int(num)
		res := Result{Method: "isPrime", IsPrime: isPrime(numInt)}
		cool, err := json.Marshal(res)
		if err != nil {
			fmt.Printf("error marshaling res")
			return
		}
		fmt.Println(cool)
		conn.Write(cool)
		conn.Write([]byte("\n"))
	}
}

// stolen from https://stackoverflow.com/questions/39950643/prime-numbers-program-gives-wrong-answer
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for ix, sqn := 3, int(math.Sqrt(float64(n))); ix <= sqn; ix += 2 {
		if n%ix == 0 {
			return false
		}
	}
	return true
}

func isIntegral(val float64) bool {
	fmt.Printf("%f\n", val)
	fmt.Printf("%f\n", float64(int(val)))
	return val == float64(int(val))
}
