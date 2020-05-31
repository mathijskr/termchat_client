package termchat

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	PORT_DEFAULT    = "31031"
	CREDENTIALS_DIR = ".config/termchat/credentials"
	DELIMITER       = "\x00"
)

func CheckErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// Send a message over a tcp connection and wait for a response.
func SendPacket(msg string, conn *net.Conn) (response string) {
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
func ReadCredentials(filename string) (username string, password string, err error) {
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
func FormatSignupPacket(username string, password string) (format string) {
	format = "signup:name=" + username + DELIMITER + "pass=" + password
	return
}

// Format the message packet according to the format the server understands.
func FormatMessagePacket(username string, password string, receiver string, body string) (format string) {
	format = "send:name=" + username + DELIMITER +
		"pass=" + password + DELIMITER +
		"receiver=" + receiver + DELIMITER +
		"body=" + body
	return
}

// Format the read contacts packet according to the format the server understands.
func FormatReadContactsPacket(username string, password string) (format string) {
	format = "contacts:name=" + username + DELIMITER +
		"pass=" + password
	return
}

func CloseConnection(conn *net.Conn) {
	SendPacket("quit:", conn)
}
