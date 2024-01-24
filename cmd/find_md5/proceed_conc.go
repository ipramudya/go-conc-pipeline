package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FileMetadata struct {
	Filepath  string
	Content   []byte
	Sum       string
	IsRenamed bool
	Err       error
}

const TotalWorkers = 100

func ProceedConcurrently() {
	/* pipeline 1: read all files */
	contents := readFiles()

	/* pipeline 2: generate checksum md5  */
	fileSum := getSum(contents)

	/* pipeline 3: rename file based on md5 sum  */
	renamedFiles := rename(fileSum)

	total := 0
	renamed := 0

	for file := range renamedFiles {
		if file.Err != nil {
			log.Printf("error renaming file %s. stack trace: %s", file.Filepath, file.Err)
		} else if file.IsRenamed {
			renamed++
		}

		total++
	}

	log.Printf("%d/%d files renamed", renamed, total)
}

func readFiles() <-chan FileMetadata {
	out := make(chan FileMetadata)

	go func() {
		err := filepath.Walk(TempDir, func(path string, info os.FileInfo, err error) error {
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
				Filepath: path,
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

	concurrentProcess(&out, func() {
		for file := range in {
			file.Sum = fmt.Sprintf("%x", md5.Sum(file.Content))
			out <- file
		}
	})

	return out
}

func rename(in <-chan FileMetadata) <-chan FileMetadata {
	out := make(chan FileMetadata)

	concurrentProcess(&out, func() {
		for file := range in {

			newPath := filepath.Join(TempDir, fmt.Sprintf("file-%s.txt", file.Sum))
			err := os.Rename(file.Filepath, newPath)
			file.IsRenamed = err == nil
			file.Err = err

			out <- file
		}
	})

	return out
}

func concurrentProcess(c *chan FileMetadata, fn func()) {
	wg := &sync.WaitGroup{}

	wg.Add(TotalWorkers)
	go func() {
		for i := 0; i < TotalWorkers; i++ {
			go func() {
				defer wg.Done()

				fn()
			}()
		}
	}()

	go func() {
		wg.Wait()
		close(*c)
	}()
}
