package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var TempDir = filepath.Join(os.TempDir(), "go-concurrency-pipeline-temp") // /tmp/go-concurrency-pipeline-temp

func main() {
	log.Println("start")
	startTime := time.Now()

	/* operations */
	// Proceed()
	ProceedConcurrently()

	duration := time.Since(startTime)
	log.Printf("done in %v seconds", duration.Seconds())
}
