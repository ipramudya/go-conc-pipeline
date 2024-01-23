package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var tempDir = filepath.Join(os.TempDir(), "go-concurrency-pipeline-temp") // /tmp/go-concurrency-pipeline-temp

type FileMetadata struct {
	FilePath  string
	Content   []byte
	Sum       string
	IsRenamed bool
}

func main() {
	log.Println("start")
	startTime := time.Now()

	/* pipeline 1 */
	contentChan := readFiles()

	/* pipeline 2 */
	fileSumChan1 := getSum(contentChan)
	fileSumChan2 := getSum(contentChan)
	fileSumChan3 := getSum(contentChan)
	fileSumChan := mergeChannels(fileSumChan1, fileSumChan2, fileSumChan3)

	/* pipeline 3 */
	renamedChan1 := rename(fileSumChan)
	renamedChan2 := rename(fileSumChan)
	renamedChan3 := rename(fileSumChan)
	renamedChan4 := rename(fileSumChan)
	renamedChan := mergeChannels(renamedChan1, renamedChan2, renamedChan3, renamedChan4)

	total := 0   // jumlah total file ditemukan
	renamed := 0 // jumlah file telah di ubah namanya
	for file := range renamedChan {
		if file.IsRenamed {
			renamed++
		}
		total++
	}

	log.Printf("%d/%d files renamed", renamed, total)

	duration := time.Since(startTime)
	log.Printf("done in %v seconds", duration.Seconds())
}

func readFiles() <-chan FileMetadata {
	out := make(chan FileMetadata)

	go func() {
		err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			buffer, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			out <- FileMetadata{
				FilePath: path,
				Content:  buffer,
			}

			return nil
		})

		if err != nil {
			log.Println("Error", err.Error())
		}

		close(out)
	}()

	return out
}

func getSum(in <-chan FileMetadata) <-chan FileMetadata {
	out := make(chan FileMetadata)

	go func() {
		for file := range in {
			file.Sum = fmt.Sprintf("%x", md5.Sum(file.Content))
			out <- file
		}
		close(out)
	}()

	return out
}

func rename(in <-chan FileMetadata) <-chan FileMetadata {
	out := make(chan FileMetadata)

	go func() {
		for file := range in {
			newPath := filepath.Join(tempDir, fmt.Sprintf("file-%s.txt", file.Sum))
			err := os.Rename(file.FilePath, newPath)
			file.IsRenamed = err == nil
			out <- file
		}

		close(out)
	}()

	return out
}

func mergeChannels(inMany ...<-chan FileMetadata) <-chan FileMetadata {
	var wg sync.WaitGroup
	out := make(chan FileMetadata)

	wg.Add(len(inMany)) // preserve waitgroup selama banyaknya argument "inMany"
	for _, eachIn := range inMany {
		go func(in <-chan FileMetadata) {

			for file := range in {
				out <- file
			}

			wg.Done()

		}(eachIn)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
