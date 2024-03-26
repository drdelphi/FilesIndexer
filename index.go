package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func index() error {
	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		_, exists := FileStrings[path]
		if err != nil {
			return err
		}

		if strings.HasSuffix(os.Args[0], path) || path == fileStringsFilename || path == configFilename || exists || info.IsDir() {
			return nil
		}

		fmt.Printf("indexing %s...\n\r", path)
		s, err := readStrings(path)
		if err == nil {
			FileStrings[path] = s
		}

		return nil
	})

	saveFileStrings()

	return err
}
