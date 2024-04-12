package main

import (
	"encoding/json"
	"log"
	"os"
)

// 108357
func main() {
	file, err := os.ReadFile("jsonfolder_combine/comments.json")
	if err != nil {
		log.Fatal(err)
	}

	data := []map[string]interface{}{}

	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	data = data[:108357]
	dataBytes, _ := json.Marshal(data)
	create, _ := os.Create("comments.json")
	create.WriteString(string(dataBytes))
	create.Close()
}
