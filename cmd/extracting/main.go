package main

import (
	"log"
	"os"
)

var (
	defaultFileURL  = "https://raw.githubusercontent.com/rit-public/HappyDB/master/happydb/data/cleaned_hm.csv"
	defaultFilePath = "./cleaned_hm.csv"
)

func main() {
	fileURL, ok := os.LookupEnv("FILE_URL")
	if !ok {
		log.Println("FILE_URL not set. Using default.")
		fileURL = defaultFileURL
	}

	filePath, ok := os.LookupEnv("FILE_PATH")
	if !ok {
		log.Println("FILE_PATH not set. Using default.")
		filePath = defaultFilePath
	}

	err := DownloadFile(filePath, fileURL)
	if err != nil {
		log.Fatal(err)
	}
}
