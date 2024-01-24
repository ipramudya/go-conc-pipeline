package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func Proceed() {
	total := 0   // jumlah total file ditemukan
	renamed := 0 // jumlah file telah di ubah namanya

	err := filepath.Walk(TempDir, func(path string, info os.FileInfo, err error) error {

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
		des := filepath.Join(TempDir, fmt.Sprintf("file-%s.txt", sum))
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
