package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

const (
	PORT_DEFAULT = "31031"
)

// Print a help message, describing the usage of the program.
func print_help() {
	fmt.Println("usage: send [-h] [-p port] server_adres username message")
}

// Send a message over a tcp connection and wait for a response.
func send_message(msg string, conn *net.Conn) (response string) {
	fmt.Fprintf(*conn, msg)
	response, err := bufio.NewReader(*conn).ReadString('\n')
	if err != nil {
		response = "error"
	}
	return
}

// Parse the command line options.
func main() {
	args := os.Args
	if len(args) < 4 {
		print_help()
		os.Exit(0)
	}

	port := flag.String("p", PORT_DEFAULT, "[port][default: "+string(PORT_DEFAULT)+"]")
	help := flag.Bool("h", false, "[help][display a help message]")

	flag.Parse()

	if *help == true {
		print_help()
		os.Exit(0)
	}

	server := os.Args[1]
	recipient := os.Args[2]
	message := os.Args[3]

	// TODO: read credentials from configuration.
	username := "joe"
	password := "secret"

	conn, err := net.Dial("tcp", server+":"+*port)
	if err != nil {
		fmt.Printf("Failed to establish a connection with %s:%s\n", server, *port)
		os.Exit(1)
	}

	fmt.Println(send_message("login", &conn))
	fmt.Println(send_message("name="+username, &conn))
	fmt.Println(send_message("pass="+password, &conn))
	fmt.Println(send_message("send", &conn))
	fmt.Println(send_message("recipient="+recipient, &conn))
	fmt.Println(send_message("message="+message, &conn))
}
