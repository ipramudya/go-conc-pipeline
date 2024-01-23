package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var tempDir = filepath.Join(os.TempDir(), "go-concurrency-pipeline-temp") // /tmp/go-concurrency-pipeline-temp

func main() {
	log.Println("start")
	startTime := time.Now()

	/* operations */
	proceed()

	duration := time.Since(startTime)
	log.Printf("done in %v seconds", duration.Seconds())
}

func proceed() {
	total := 0   // jumlah total file ditemukan
	renamed := 0 // jumlah file telah di ubah namanya

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		total++

		buffer, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		/* got md5 sum */
		sum := fmt.Sprintf("%x", md5.Sum(buffer))

		/* rename file */
		des := filepath.Join(tempDir, fmt.Sprintf("file-%s.txt", sum))
		err = os.Rename(path, des)
		if err != nil {
			return err
		}

		renamed++

		return nil
	})

	if err != nil {
		log.Println("ERROR:", err.Error())
	}

	log.Printf("%d/%d files renamed", renamed, total)
}
