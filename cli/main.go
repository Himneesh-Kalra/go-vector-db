package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Himneesh-Kalra/go-vec-cli/config"
	"github.com/c-bata/go-prompt"
)

var (
	conn       net.Conn
	shouldExit = false
)

func main() {
	fmt.Println("Welcome to Vec-CLI")

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("config error: %v", err)
	}
	fmt.Println("*****Config Loaded*****")

	tcpAddr := config.CLIConfig.TCPAddr
	var err error
	conn, err = net.Dial("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Failed to connect ", err)
		os.Exit(1)
	}
	defer conn.Close()

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("vecDB> "),
		prompt.OptionTitle("VecDB CLI"),
	)

	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nRecieved interrupt. Exiting ...")
		shouldExit = true
		cancel()
		conn.Close()
		os.Exit(0)
	}()

	fmt.Println("Connected to VecDB CLI Server at: ", tcpAddr)
	go listenForServer(ctx)

	p.Run()
}

func executor(input string) {
	if conn == nil {
		fmt.Println("Not connected. Exiting...")
		os.Exit(1)
	}

	_, err := conn.Write([]byte(input + "\n"))
	if err != nil {
		fmt.Println("Error sending command : ", err)
	}

	if strings.ToUpper(input) == "EXIT" {
		return
	}
}

func listenForServer(ctx context.Context) {
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if shouldExit {
				return
			}
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Connection Lost")
				conn.Close()
				conn = nil
				os.Exit(1)
			}
			fmt.Printf("Server: %s", message)
		}
	}

}

func completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "INSERT", Description: "Insert a vector into a table"},
		{Text: "GETALL", Description: "Get all the vectors in a specific table"},
		{Text: "GETK", Description: "Search TopK elements in a table"},
		{Text: "DELETE", Description: "Delete a vector from a specific table"},
		{Text: "PING", Description: "Ping the server"},
		{Text: "EXIT", Description: "Exit the CLI"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
