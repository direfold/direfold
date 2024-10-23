package main

import (
	"errors"
	"flag"
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
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "verbose logs")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `usage: %s [flags] [path]

default [path] is "."

flags:
`, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	var path string
	if args := flag.Args(); len(args) > 0 {
		path = args[0]
	}

	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

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
	if err := filepath.WalkDir(path, walkFunc); err != nil {
		log.Printf("can't walk dir: %v", err)
		return
	}

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

	if verbose {
		fmt.Println()
		log.Printf("searched for %s :)", time.Since(start))
	}
}
