package main

import (
	"bufio"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

type Work struct {
	file  string
	regex *regexp.Regexp
}

func main() {
	start := time.Now()
	defer func() {
		dur := time.Since(start)
		slog.Info("Operation finished",
			slog.Duration("execution_time", dur),
		)
	}()
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <directory> <regex>")
	}

	workCh := make(chan Work, 100)

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)

		go func() {
			defer wg.Done()

			Worker(workCh)
		}()
	}

	re, err := regexp.Compile(os.Args[2])
	if err != nil {
		log.Fatalf("compile regex %s", os.Args[2])
	}

	var senderWg sync.WaitGroup

	senderWg.Add(1)
	go func() {
		defer senderWg.Done()

		filepath.Walk(os.Args[1], func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				slog.Warn("unable to walk path", "path", path, "error", err)
				return nil
			}

			if info.IsDir() {
				return nil
			}

			workCh <- Work{
				file:  path,
				regex: re,
			}
			return nil
		})
	}()

	senderWg.Wait()
	close(workCh)
	wg.Wait()
}

func Worker(workCh chan Work) {
	for work := range workCh {
		processWork(work.file, work.regex)
	}
}

func processWork(fn string, re *regexp.Regexp) {
	f, err := os.Open(fn)
	if err != nil {
		slog.Error("unable to open file", "file", fn, "error", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	scn := bufio.NewScanner(f)
	lineNum := 1
	for scn.Scan() {
		result := re.Find(scn.Bytes())
		if len(result) > 0 {
			slog.Info("Match found", "file", fn, "line", lineNum, "match", string(result))
		}

		lineNum++
	}
}
