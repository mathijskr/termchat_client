package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
	"strings"

	tc "github.com/mathijskr/termchat_client/termchat"
)

// Print a help message, describing the usage of the program.
func printHelp() {
	fmt.Println("usage: read [-h help] [-p port] [-a server_adres]")
}

// Format the read chat packet according to the format the server understands.
func FormatReadChatPacket(username string, password string, contact string, timestamp string) (format string) {
	format = "read:name=" + username + tc.DELIMITER +
		"pass=" + password + tc.DELIMITER +
		"contact=" + contact + tc.DELIMITER +
		"timestamp=" + timestamp
	return
}

// Remove the trailing newline and null character.
func removeTrailing(s string) string {
	if len(s) > 2 {
		return s[:len(s)-2]
	}
	return ""
}

// Find the timestamp of the last saved message.
func lastMessageTimestamp(chatFile *os.File) string {
	// Find the last locally saved message.
	scanner := bufio.NewScanner(chatFile)
	lastMessage := ""
	for scanner.Scan() {
		lastMessage = scanner.Text()
	}

	// Extract the timestamp from the message.
	timestampIndex := 0
	for index, c := range lastMessage {
		if c == ':' {
			break
		}
		timestampIndex = index
	}
	if timestampIndex > 0 {
		return lastMessage[0 : timestampIndex+1]
	}
	return "0"
}

// Update the local chat history for a contact.
func updateChat(contact string, username string, password string, chatDir string, conn *net.Conn) {
	// Create the chat directory if it doesn't exist.
	err := os.MkdirAll(chatDir, 0755)
	if !tc.CheckErr(err) {
		fmt.Println("Cannot create directory: ", chatDir)
		return
	}

	chatFile, err := os.OpenFile(chatDir+"/"+contact, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if !tc.CheckErr(err) {
		return
	}

	defer chatFile.Close()

	timestamp := lastMessageTimestamp(chatFile)

	readChatPacket := FormatReadChatPacket(username, password, contact, timestamp)
	chatPacket := tc.SendPacket(readChatPacket, conn)
	chat := strings.Split(removeTrailing(chatPacket), "\x00")

	for _, field := range chat {
		if strings.HasPrefix(field, "timestamp=") {
			_, err = chatFile.WriteString(field[len("timestamp="):len(field)] + ":")
			if !tc.CheckErr(err) {
				return
			}
		}
		if strings.HasPrefix(field, "sender=") {
			_, err = chatFile.WriteString(field[len("sender="):len(field)] + ":")
			if !tc.CheckErr(err) {
				return
			}
		}
		if strings.HasPrefix(field, "body=") {
			_, err = chatFile.WriteString(field[len("body="):len(field)] + "\n")
			if !tc.CheckErr(err) {
				return
			}
		}
	}

	return
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
		fmt.Println("Please store your credentials in ~/.config/termchat/credentials/<server_adres>.conf")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", *server+":"+*portFlag)
	if err != nil {
		fmt.Printf("Failed to establish a connection with %s:%s\n", *server, *portFlag)
		os.Exit(1)
	}

	readContactsPacket := tc.FormatReadContactsPacket(username, password)
	contactsPacket := tc.SendPacket(readContactsPacket, &conn)
	contacts := strings.Split(removeTrailing(contactsPacket), "\x00")

	// Read chat for every contact.
	for _, contact := range contacts {
		if contact != "" {
			updateChat(contact, username, password, "chats", &conn)
		}
	}

	tc.CloseConnection(&conn)
}
