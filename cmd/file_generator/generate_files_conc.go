package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FileMetadata struct {
	Index       int
	Filename    string
	WorkerIndex int
	Err         error
}

func GenerateFilesConcurently() {
	os.RemoveAll(TempDir)
	os.MkdirAll(TempDir, os.ModePerm)

	/* pipeline 1: job distribution */
	fileIndexChan := generateFileIndexes()

	/* pipeline 2: creating files */
	createFilesWorker := 100
	fileResultChan := createFiles(fileIndexChan, createFilesWorker)

	total := 0
	created := 0

	for file := range fileResultChan {
		if file.Err != nil {
			log.Printf("error creating file %s. stack trace: %s", file.Filename, file.Err)
		} else {
			created++
		}
		total++
	}

	log.Printf("%d/%d of total files created", created, total)
}

func generateFileIndexes() <-chan FileMetadata {
	out := make(chan FileMetadata)

	go func() {
		defer close(out)

		for i := 0; i < TotalFiles; i++ {
			file := FileMetadata{
				Index:    i,
				Filename: fmt.Sprintf("file-%d.txt", i),
			}

			out <- file
		}

	}()

	return out
}

func createFiles(in <-chan FileMetadata, totalWorkers int) <-chan FileMetadata {
	out := make(chan FileMetadata)

	wg := &sync.WaitGroup{}

	wg.Add(totalWorkers)
	go func() {
		/* dispatch n workers */
		for workerIndex := 0; workerIndex < totalWorkers; workerIndex++ {

			/* listen to channel in for incoming jobs */
			go func(i int) {

				/* every each job done, subtract the wait group amount */
				defer wg.Done()

				/* each channel in's job is represented by file */
				for file := range in {
					path := filepath.Join(TempDir, file.Filename)
					content := RandomString(ContentLength)
					err := os.WriteFile(path, []byte(content), os.ModePerm)

					/* construct job's result, send it through channel out */
					out <- FileMetadata{
						Filename:    file.Filename,
						WorkerIndex: i,
						Err:         err,
					}
				}

			}(workerIndex)

		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
