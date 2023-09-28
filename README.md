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


On a Chromebook, share with Linux your main folders on Google Drive. Afterwards, 
you can go to the mount location and see them. For example:

```
$ pwd
/mnt/chromeos/GoogleDrive/MyDrive
$ pwd
/mnt/chromeos/GoogleDrive/MyDrive
$ ls
Backup  Documents  Music  Pictures  Projects  Videos
$ 
```

Next build and run FileReporter:

```
$ go build FileReporter.go
$ cp FileReporter ~/bin # somewhere on your path!
$ cd /mnt/chromeos/GoogleDrive/MyDrive
$ FileReporter . ~/gdrive.csv
```

You can a `tail` follow to track progress:

```
tail -f grive.csv
```
