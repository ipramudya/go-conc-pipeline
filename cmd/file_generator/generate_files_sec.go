package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GenerateFiles() {
	os.RemoveAll(TempDir)
	os.MkdirAll(TempDir, os.ModePerm)

	for i := 0; i < TotalFiles; i++ {
		filename := filepath.Join(TempDir, fmt.Sprintf("file-%d.txt", i))
		content := RandomString(ContentLength)

		err := os.WriteFile(filename, []byte(content), os.ModePerm)
		if err != nil {
			log.Println("error writing file", filename)
		}

		/* log every hundred */
		// if i%100 == 0 && i > 0 {
		// 	log.Println(i, "files created")
		// }
	}

	log.Printf("%d of total files created", TotalFiles)
}
