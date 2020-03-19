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

// Read the credentials from a user.
// Expects username to be the first string in a file
// and password the second string.
func read_credentials(filename string) (username string, password string, err error) {
	username = ""
	password = ""

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Fail")
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

func main() {
	args := os.Args
	if len(args) < 4 {
		print_help()
		os.Exit(0)
	}

	// Parse the command line options.
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

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to locate home directory")
		os.Exit(1)
	}

	username, password, err := read_credentials(usr.HomeDir + "/" + CREDENTIALS_DIR + "/" + server + ".conf")
	if err != nil {
		fmt.Printf("Please store your credentials in ~/.termchat/credentials/<server_adres>.conf")
		os.Exit(1)
	}

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
