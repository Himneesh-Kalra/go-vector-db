package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/Himneesh-Kalra/go-vector-db/api"
	"github.com/Himneesh-Kalra/go-vector-db/config"
	"github.com/Himneesh-Kalra/go-vector-db/db"
	"github.com/Himneesh-Kalra/go-vector-db/tcp"

	// "github.com/Himneesh-Kalra/go-vector-db/search"
	"github.com/Himneesh-Kalra/go-vector-db/storage"
)

func main() {

	var wg sync.WaitGroup

	//load via viper
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Config error: %v", err)
	}
	fmt.Println("****** Config Loaded ******")
	//init storage
	storage.DataDir = config.AppConfig.DataDir
	fmt.Printf("****** Storage directory Configured :%v ******\n", storage.DataDir)

	if err := storage.LoadAll(); err != nil {
		log.Fatal("Failed to load data from storage:", err)
	}
	fmt.Println("****** Data Successfully loaded from disk ******")

	//Algo selection

	switch config.AppConfig.SearchAlgorithm {
	case "brute":
		db.UseBrute()
		fmt.Println("******* Using Brute Search ********")

	case "lsh":
		db.UseLSH(10)
		fmt.Println("******* Using LSH Search ********")
	default:
		log.Fatalf("Unknown algorithm: %s", config.AppConfig.SearchAlgorithm)
	}
	addr := fmt.Sprintf(":%d", config.AppConfig.ServerPort)
	tcpAddr := config.AppConfig.TCPAddr

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting TCP listener on: ", tcpAddr)
		tcp.StartTCPServer(tcpAddr)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app := api.SetupRoutes()
		log.Printf("****** Listening on %s ******", addr)
		log.Fatal(app.Listen(addr))
	}()

	wg.Wait()

}
