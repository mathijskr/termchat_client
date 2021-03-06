package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"

	tc "github.com/mathijskr/termchat_client/termchat"
)

// Print a help message, describing the usage of the program.
func printHelp() {
	fmt.Println("usage: send [-h help] [-s signup] [-p port] [-a server_adres] [-r receiver] [-m message]")
}

func main() {
	// Parse the command line options.
	portFlag := flag.String(
		"p", tc.PORT_DEFAULT, "[port][default: "+string(tc.PORT_DEFAULT)+"]",
	)
	helpFlag := flag.Bool(
		"h", false, "[help][display a help message]",
	)
	signupFlag := flag.Bool(
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
		((*receiver == "" ||
			*body == "") &&
			!*signupFlag) ||
		*helpFlag {
		printHelp()
		os.Exit(0)
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Println("Failed to locate home directory")
		os.Exit(1)
	}

	username, password, err := tc.ReadCredentials(usr.HomeDir + "/" + tc.CREDENTIALS_DIR + "/" + *server + ".conf")
	if err != nil {
		fmt.Println("Please store your credentials in ~/.config/termchat/credentials/<server_adres>.conf")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", *server+":"+*portFlag)
	if err != nil {
		fmt.Printf("Failed to establish a connection with %s:%s\n", *server, *portFlag)
		os.Exit(1)
	}

	switch {
	case *signupFlag:
		fmt.Println("Signing up")
		signup := tc.FormatSignupPacket(username, password)
		fmt.Printf(tc.SendPacket(signup, &conn))
	default:
		fmt.Println("Sending message")
		message := tc.FormatMessagePacket(username, password, *receiver, *body)
		fmt.Printf(tc.SendPacket(message, &conn))
	}

	tc.CloseConnection(&conn)
}
