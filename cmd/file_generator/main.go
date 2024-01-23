package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const totalFile = 3000
const contentLength = 5000

var tempDir = filepath.Join(os.TempDir(), "go-concurrency-pipeline-temp") // /tmp/go-concurrency-pipeline-temp

func main() {
	log.Println("start")
	startTime := time.Now()

	generateFiles()

	duration := time.Since(startTime)
	log.Printf("done in %v seconds", duration.Seconds())
}

func generateFiles() {
	os.RemoveAll(tempDir)
	os.MkdirAll(tempDir, os.ModePerm)

	for i := 0; i < totalFile; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("file-%d.txt", i))
		content := randomString(contentLength)

		err := os.WriteFile(filename, []byte(content), os.ModePerm)
		if err != nil {
			log.Println("error writing file", filename)
		}

		/* log every hundred */
		if i%100 == 0 && i > 0 {
			log.Println(i, "files created")
		}
	}

	log.Printf("%d of total files created", totalFile)
}

func randomString(length int) string {
	randomizer := rand.New(rand.NewSource(time.Now().Unix()))
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)

	for i := range b {
		b[i] = letters[randomizer.Intn(len(letters))]
	}

	return string(b)
}
