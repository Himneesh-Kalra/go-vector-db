package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func StartTCPServer(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v \n", err)
	}
	fmt.Printf("TCP server listening on %s\n", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed")
			return
		}
		command = strings.TrimSpace(command)
		args := strings.SplitN(command, " ", 2)

		switch strings.ToUpper(args[0]) {
		case "INSERT":
			InsertHandler(conn, args)

		case "GETALL":
			GetAllHandler(conn, args)

		case "GETK":
			GetKHandler(conn, args)

		case "DELETE":
			DeleteHandler(conn, args)

		case "PING":
			fmt.Fprintln(conn, "PONG")

		case "EXIT":
			fmt.Fprintln(conn, "GoodBye!! ...")
			return

		default:
			fmt.Fprintln(conn, "Unknown Command")
		}
	}
}
