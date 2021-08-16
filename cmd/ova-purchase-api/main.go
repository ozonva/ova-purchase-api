package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Hey! This is ova purchase api service 💸")
	LoadConfig()
}

type Config struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func loadJson(filename string) {
	config := Config{}

	fmt.Println("Opening file...")
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		fmt.Println("Closing file....")
		err := file.Close()
		if err != nil {
			log.Fatal("Oh no, can't close file!")
		}
	}()

	reader := json.NewDecoder(file)
	if err := reader.Decode(&config); err != nil {
		fmt.Printf("Error parsing json! %s", err)
	}
	fmt.Printf("Config(%v)\n", config)
}

func LoadConfig() {
	for i := 0; i < 10; i++ {
		loadJson("config.json")
	}
}
