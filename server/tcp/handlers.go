package tcp

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Himneesh-Kalra/go-vector-db/db"
	"github.com/Himneesh-Kalra/go-vector-db/models"
	"github.com/Himneesh-Kalra/go-vector-db/storage"
)

// CLI GETALL Vectors
func GetAllHandler(conn net.Conn, args []string) {
	if len(args) != 2 {
		fmt.Fprintln(conn, "Usage: GETALL <tablename>")
		return
	}
	tableName := args[1]
	data, ok := storage.GetTable(tableName)
	if !ok {
		fmt.Fprintln(conn, "Table does not exist")
		return
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintln(conn, "Failed to serialize data")
		return
	}

	fmt.Fprintf(conn, "Data from %s:\n%s\n", tableName, string(jsonData))
}

// CLI Insert Vector
func InsertHandler(conn net.Conn, args []string) {
	fmt.Println("args: ", args)
	fmt.Println("lengts of args : ", len(args))

	if len(args) < 2 {
		fmt.Fprintln(conn, "Usage INSERT <tablename> <jsonpayload>")
		return
	}

	rawCommand := args[1]
	fmt.Println("Raw command: ", rawCommand)

	spaceIndex := strings.Index(rawCommand, " ")
	fmt.Println("first space occurence: ", spaceIndex)
	if spaceIndex == -1 {
		fmt.Fprintln(conn, "Invalid Format. Json payload is missing...")
		return
	}
	tableName := rawCommand[:spaceIndex]
	payload := string(rawCommand[spaceIndex+1:])

	if strings.HasPrefix(payload, "'") && strings.HasSuffix(payload, "'") {
		payload = strings.Trim(payload, "'")
	}

	fmt.Printf("Received payload : %s\n", payload)

	var vector models.Vector
	if err := json.Unmarshal([]byte(payload), &vector); err != nil {
		fmt.Fprintln(conn, "Invalid JSON format")
		return
	}

	storage.InsertVector(tableName, vector)
	fmt.Fprintf(conn, "Vector inserted into table %s\n", tableName)
}

// CLI GETK vectors
func GetKHandler(conn net.Conn, args []string) {
	fmt.Println("args: ", args)
	fmt.Println("Length of args: ", len(args))

	if len(args) < 2 {
		fmt.Fprintln(conn, "Usage GETK <tablename> <k> <jsonpayload>")
		return
	}

	rawCommand := args[1]
	fmt.Println("rawcommand: ", rawCommand)

	firstSpaceIndex := strings.Index(rawCommand, " ")
	fmt.Println("First space occurence: ", firstSpaceIndex)
	if firstSpaceIndex == -1 {
		fmt.Fprintln(conn, "Invalid format, json payload is missing...")
		return
	}
	tableName := rawCommand[:firstSpaceIndex]
	fmt.Println("Table name: ", tableName)

	afterTableName := rawCommand[firstSpaceIndex+1:]

	secondSpaceIndex := strings.Index(afterTableName, " ")
	if secondSpaceIndex == -1 {
		fmt.Fprintln(conn, "invalid format, k value is missing ")
		return
	}
	secondSpaceIndex += firstSpaceIndex + 1
	fmt.Println("second space occurence: ", secondSpaceIndex)

	rawK := rawCommand[firstSpaceIndex+1 : secondSpaceIndex]
	k, err := strconv.Atoi(rawK)
	if err != nil {
		fmt.Println("Error. Could not parse K ", err)
		return
	}
	fmt.Println("Value of K: ", k)

	payload := string(rawCommand[secondSpaceIndex+1:])

	if strings.HasPrefix(payload, "'") && strings.HasSuffix(payload, "'") {
		payload = strings.Trim(payload, "'")
	}

	fmt.Println("payload recieved: ", payload)

	var query []float32
	if err := json.Unmarshal([]byte(payload), &query); err != nil {
		fmt.Fprintln(conn, "invalid json format")
		return
	}
	results := db.SearchTopK(query, k, tableName)
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Fprintln(conn, "Failed to serialize data")
		return
	}
	fmt.Fprintf(conn, "TOPK Results from Table: %s\n%s\n", tableName, string(jsonData))

}

// CLI Delete Vector
func DeleteHandler(conn net.Conn, args []string) {
	if len(args) < 2 {
		fmt.Fprintln(conn, "Usage DELETE <tablename> <id>")
	}

	rawCommand := args[1]
	fmt.Println("raw command: ", rawCommand)

	spaceIndex := strings.Index(rawCommand, " ")
	fmt.Println("first space occurence at index: ", spaceIndex)

	if spaceIndex == -1 {
		fmt.Fprintln(conn, "invalid format, id is missing")
		return
	}

	tableName := rawCommand[:spaceIndex]
	fmt.Println("Table name: ", tableName)

	vectorID := rawCommand[spaceIndex+1:]

	deleted := storage.DeleteVector(tableName, vectorID)
	if !deleted {
		fmt.Fprintln(conn, "Could not delete vector")
		return
	}
	fmt.Fprintf(conn, "Vector: %s deleted from table: %s  Successfully...", vectorID, tableName)
}
