package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"

	tc "github.com/mathijskr/termchat/Client/termchat"
)

// Print a help message, describing the usage of the program.
func printHelp() {
	fmt.Println("usage: read [-h help] [-p port] [-a server_adres]")
}

func main() {
	// Parse the command line options.
	portFlag := flag.String(
		"p", tc.PORT_DEFAULT, "[port][default: "+string(tc.PORT_DEFAULT)+"]",
	)
	helpFlag := flag.Bool(
		"h", false, "[help][display a help message]",
	)
	server := flag.String(
		"a", "", "[server_adres]",
	)
	flag.Parse()

	if *server == "" ||
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
		fmt.Println("Please store your credentials in ~/.termchat/credentials/<server_adres>.conf")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", *server+":"+*portFlag)
	if err != nil {
		fmt.Printf("Failed to establish a connection with %s:%s\n", *server, *portFlag)
		os.Exit(1)
	}

	fmt.Println("Sending read contacts request")
	read := tc.FormatReadContactsPacket(username, password)
	fmt.Printf(tc.SendPacket(read, &conn))
}
