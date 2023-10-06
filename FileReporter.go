package main

// FileReporter generates a CSV report of all files starting in
// the specified directory. This is a concurrent implementation
// based on the example from gopl.io in chapter 8.
// It uses a concurrency-limiting counting semaphore
// to avoid opening too many files at once.

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// PathfileInfo captures both file info and the path to the file
type PathfileInfo struct {
	os.FileInfo
	Path string
}

var numthreads = flag.Int("n", 20, "Number of go routines")
var hidden = flag.Bool("h", false, "Process hidden directories")
var quiet = flag.Bool("q", true, "Suppress directory permission errors")
var help = flag.Bool("help", false, "Show help message")

func main() {
	flag.Parse()

	if *help {
		fmt.Println("Usage: FileReporter dirname output.csv")
		fmt.Println("Note: 'hidden' means a prefix of a 'dot' (default false)")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if len(flag.Args()) < 2 {
		fmt.Println("Usage: FileReporter dirname output.csv")
		fmt.Println("Note: 'hidden' means a prefix of a 'dot' (default false)")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// open output file
	fo, foerr := os.Create(flag.Args()[1])
	if foerr != nil {
		fmt.Fprintf(os.Stderr, "Error on os.Create():\n%v\n", foerr)
		os.Exit(-1)
	}
	defer fo.Close()
	w := csv.NewWriter(fo)
	err := w.Write([]string{"Size", "Filename", "Ext", "Directory", "FullPath"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on csv.Write():\n%v\n", err)
		os.Exit(-1)
	}

	files := make(chan PathfileInfo)
	var wg sync.WaitGroup
	wg.Add(1)
	go walkDir(flag.Args()[0], &wg, files)

	// this needed to actually close the channel when all the
	// go routines are done. Otherwise, this error will be
	// returned when everything is done:
	// fatal error: all goroutines are asleep - deadlock!
	go func() {
		wg.Wait()
		close(files)
	}()

	for {
		finfo, ok := <-files
		if !ok {
			break // files channel closed
		}
		_name := finfo.Name()
		_size := finfo.Size()
		_path := finfo.Path
		err := w.Write([]string{
			fmt.Sprintf("%v", _size),
			_name,
			filepath.Ext(_name),
			_path,
			filepath.Join(finfo.Path, finfo.Name())})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on csv.Write():\n%v\n", err)
			os.Exit(-1)
		}
	}
	w.Flush()

}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(dir string, wg *sync.WaitGroup, f chan<- PathfileInfo) {
	defer wg.Done()

	for _, entry := range dirents(dir) {
		isADir := entry.IsDir()
		if isADir {
			if strings.HasPrefix(entry.Name(), ".") {
				if !*hidden {
					continue
				}
			}

			wg.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, wg, f)
		} else {
			_finfo, err := entry.Info()
			if err == nil {
				pfinfo := PathfileInfo{_finfo, dir}
				f <- pfinfo
			} else {
				fmt.Fprintf(os.Stderr, "Error on ReadDir():\n%v\n", err)
			}
		}
	}
}

var sema = make(chan struct{}, *numthreads)

// dirents returns the entries of directory dir.
func dirents(dir string) []os.DirEntry {
	sema <- struct{}{}        // acquire token, max 20
	defer func() { <-sema }() // release token

	entries, err := os.ReadDir(dir)
	/*
			have to be careful with this test. May need to leave
			commented out, since if directories are blinking
			in and out of existence, we can't accept this
			as a hard failure!
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error on ReadDir():\n%v\n", err)
			os.Exit(-1)
		}
	*/
	if err != nil {
		if os.IsPermission(err) && *quiet {
			return nil
		}
		fmt.Fprintf(os.Stderr, "Error on ReadDir():\n%v\n", err)
		return nil
	}
	return entries
}
