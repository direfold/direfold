package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [path to directory]\n", os.Args[0])
		return
	}

	path := filepath.Clean(os.Args[1])

	var dirs int64
	var files int64

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type().IsDir() {
			dirs++
		}
		if d.Type().IsRegular() {
			files++
		}

		return nil
	}

	err := filepath.WalkDir(path, walkFunc)
	if err != nil {
		log.Printf("can't walk dir: %v", err)
		return
	}

	output := []struct {
		text    string
		value   any
		leftpad int
	}{
		{"", "", 0},
		{"direfold: ", path, 0},
		{"---", "", 0},
		{"size: ", humanize.Bytes(0), 10},
		{"files: ", humanize.Comma(files), 13},
		{"directories: ", humanize.Comma(dirs), 0},
	}

	for _, line := range output {
		// fmt.Printf("%s%*s\n", line.text, line.leftpad, line.value)
		fmt.Printf("%s%s\n", line.text, line.value)
	}

	fmt.Println()
}
