# FileReporter
Produces a report (CSV) of files.

Usage:
```
$ go run FileReporter.go -help
Usage: FileReporter dirname output.csv
Note: 'hidden' means a prefix of a 'dot' (default false)
  -h	Process hidden directories
  -help
    	Show help message
  -n int
    	Number of go routines (default 20)
  -q	Suppress directory permission errors (default true)
```