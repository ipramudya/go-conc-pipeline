package main

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const TotalFiles = 10000
const ContentLength = 5000

var TempDir = filepath.Join(os.TempDir(), "go-concurrency-pipeline-temp") // /tmp/go-concurrency-pipeline-temp

func main() {
	log.Println("start")
	startTime := time.Now()

	// GenerateFiles()
	GenerateFilesConcurently()

	duration := time.Since(startTime)
	log.Printf("done in %v seconds", duration.Seconds())
}

func RandomString(length int) string {
	randomizer := rand.New(rand.NewSource(time.Now().Unix()))
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)

	for i := range b {
		b[i] = letters[randomizer.Intn(len(letters))]
	}

	return string(b)
}
