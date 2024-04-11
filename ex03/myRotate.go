package main

import (
    "archive/tar"
    "compress/gzip"
    "flag"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sync"
    "time"
)

func createTarGz(source, target string) error {
    targetFile, err := os.Create(target)
    if err != nil {
        return err
    }
    defer targetFile.Close()

    gzipWriter := gzip.NewWriter(targetFile)
    defer gzipWriter.Close()

    tarWriter := tar.NewWriter(gzipWriter)
    defer tarWriter.Close()

    sourceFile, err := os.Open(source)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    fileInfo, err := sourceFile.Stat()
    if err != nil {
        return err
    }

    header := &tar.Header{
        Name:    filepath.Base(source),
        Mode:    int64(fileInfo.Mode()),
        ModTime: fileInfo.ModTime(),
        Size:    fileInfo.Size(),
    }

    if err := tarWriter.WriteHeader(header); err != nil {
        return err
    }

    if _, err := io.Copy(tarWriter, sourceFile); err != nil {
        return err
    }

    return nil
}

func rotateLogFile(logFile, archiveDir string) error {
    timestamp := time.Now().Unix()

    archiveFile := fmt.Sprintf("%s_%d.tar.gz", logFile, timestamp)
    archivePath := filepath.Join(archiveDir, filepath.Base(archiveFile))

    if err := os.MkdirAll(archiveDir, 0755); err != nil {
        return err
    }

    if err := createTarGz(logFile, archivePath); err != nil {
        return err
    }

    return nil
}

func main() {
    archiveDir := flag.String("a", "", "archive directory")
    flag.Parse()

    if *archiveDir == "" {
        fmt.Println("Error: Archive directory not provided.")
        os.Exit(1)
    }

    logFiles := flag.Args()

    var wg sync.WaitGroup

    for _, logFile := range logFiles {
        wg.Add(1)
        go func(logFile string) {
            defer wg.Done()
            if err := rotateLogFile(logFile, *archiveDir); err != nil {
                fmt.Printf("Error rotating %s: %v\n", logFile, err)
            } else {
                fmt.Printf("Rotated %s successfully\n", logFile)
            }
        }(logFile)
    }

    wg.Wait()
}
