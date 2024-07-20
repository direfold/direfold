package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

var ignore = []string{"/proc/", "/sys/", "/run/"}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [path to directory]\n", os.Args[0])
		return
	}

	path := filepath.Clean(os.Args[1])

	var dirs int64
	var files int64
	var sizes int64

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		for _, prefix := range ignore {
			if strings.HasPrefix(path, prefix) {
				return fs.SkipDir
			}
		}

		if errors.Is(err, fs.ErrPermission) {
			log.Println(err)
		} else if err != nil {
			return fmt.Errorf("walk func err: %w", err)
		}

		if d.Type().IsDir() {
			dirs++
		}
		if d.Type().IsRegular() {
			files++

			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("can't get fileinfo for path %q: %w", path, err)
			}

			sizes += info.Size()
		}

		return nil
	}

	start := time.Now()

	err := filepath.WalkDir(path, walkFunc)
	if err != nil {
		log.Printf("can't walk dir: %v", err)
		return
	}

	log.Printf("walk took %s", time.Since(start))

	output := []struct {
		text  string
		value string
	}{
		{"", ""},
		{"direfold: ", path},
		{"---", ""},
		{"size:    ", humanize.Bytes(uint64(sizes))},
		{"files:   ", humanize.Comma(files)},
		{"folders: ", humanize.Comma(dirs)},
	}

	for _, line := range output {
		fmt.Printf("%s%s\n", line.text, line.value)
	}

	fmt.Println()
}
