package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	PORT_DEFAULT = 31031
)

func print_help() {
	fmt.Println("usage: send [-h] [-p port] server_adres username message")
}

func main() {
	args := os.Args
	if len(args) < 4 {
		print_help()
		os.Exit(0)
	}

	port := flag.Int("p", PORT_DEFAULT, "[port][default: "+string(PORT_DEFAULT)+"]")
	help := flag.Bool("h", false, "[help][display a help message]")

	flag.Parse()

	if *help == true {
		print_help()
		os.Exit(0)
	}

	recipient := os.Args[1]
	server := os.Args[2]
	message := os.Args[3]

	fmt.Println(*port, recipient, server, message)
}
