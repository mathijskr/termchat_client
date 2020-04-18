package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
)

const (
	PORT_DEFAULT    = "31031"
	CREDENTIALS_DIR = ".termchat/credentials"
	DELIMITER       = "\x00"
)

// Print a help message, describing the usage of the program.
func print_help() {
	fmt.Println("usage: send [-h help] [-s signup] [-p port] [-a server_adres] [-r receiver] [-m message]")
}

// Send a message over a tcp connection and wait for a response.
func send_packet(msg string, conn *net.Conn) (response string) {
	fmt.Fprintf(*conn, msg)
	response, err := bufio.NewReader(*conn).ReadString('\n')
	if err != nil {
		response = "error"
	}
	return
}

// Read the credentials from a user.
// Expects username to be the first string in a file
// and password the second string.
func read_credentials(filename string) (username string, password string, err error) {
	username = ""
	password = ""

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	username = scanner.Text()
	scanner.Scan()
	password = scanner.Text()
	return
}

// Format the signup packet according to the format the server understands.
func format_signup_packet(username string, password string) (format string) {
	format = "signup:name=" + username + DELIMITER + "pass=" + password
	return
}

// Format the message packet according to the format the server understands.
func format_message_packet(username string, password string, receiver string, body string) (format string) {
	format = "send:name=" + username + DELIMITER +
		"pass=" + password + DELIMITER +
		"receiver=" + receiver + DELIMITER +
		"body=" + body
	return
}

// Format the read packet according to the format the server understands.
func format_read_packet(username string, password string, contact string) (format string) {
	format = "read:name=" + username + DELIMITER +
		"pass=" + password + DELIMITER +
		"contact=" + contact
	return
}

func main() {
	// Parse the command line options.
	port_flag := flag.String(
		"p", PORT_DEFAULT, "[port][default: "+string(PORT_DEFAULT)+"]",
	)
	help_flag := flag.Bool(
		"h", false, "[help][display a help message]",
	)
	signup_flag := flag.Bool(
		"s", false, "[signup][create an account at the server you are connecting to]",
	)
	server := flag.String(
		"a", "", "[server_adres]",
	)
	receiver := flag.String(
		"r", "", "[receiver]",
	)
	body := flag.String(
		"m", "", "[message]",
	)
	flag.Parse()

	if *server == "" ||
		*receiver == "" ||
		*body == "" {
		print_help()
		os.Exit(0)
	}

	if *help_flag {
		print_help()
		os.Exit(0)
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Failed to locate home directory")
		os.Exit(1)
	}

	username, password, err := read_credentials(usr.HomeDir + "/" + CREDENTIALS_DIR + "/" + *server + ".conf")
	if err != nil {
		fmt.Println("Please store your credentials in ~/.termchat/credentials/<server_adres>.conf")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", *server+":"+*port_flag)
	if err != nil {
		fmt.Printf("Failed to establish a connection with %s:%s\n", *server, *port_flag)
		os.Exit(1)
	}

	switch {
	case *signup_flag:
		fmt.Println("Signing up")
		signup := format_signup_packet(username, password)
		fmt.Printf(send_packet(signup, &conn))
	//case *read_flag:
	//	fmt.Println("Sending read request")
	//	contact := "Banaan"
	//	read := format_read_packet(username, password, contact)
	//	fmt.Printf(send_packet(read, &conn))
	default:
		fmt.Println("Sending message")
		message := format_message_packet(username, password, *receiver, *body)
		fmt.Printf(send_packet(message, &conn))
	}
}
