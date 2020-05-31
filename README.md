# termchat_client

A simple go client for connecting to a termchat server.

Usage:
* Put your credentials for the servers you want to connect to in: ~/.config/termchat/credentials/<servername>.conf
* Execute "go run cmd/tc_send/main.go -a <servername> -r <receiver> -m <message>" to send a message
* Execute "go run cmd/tc_read/main.go -a <servername>" to read your inbox, chats are saved in ./chats/<contact>
