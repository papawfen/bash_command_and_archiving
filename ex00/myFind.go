package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: myFind [-f] [-d] [-sl] [-ext extension] path")
        os.Exit(1)
    }

    // Parse command line arguments
    var findFiles, findDirs, findSymlinks bool
    var fileExt string
    path := os.Args[len(os.Args) - 1]
    for _, arg := range os.Args[1 : len(os.Args) - 1] {
        switch arg {
        case "-f":
            findFiles = true
        case "-d":
            findDirs = true
        case "-sl":
            findSymlinks = true
        case "-ext":
            if len(os.Args) < 4 {
                fmt.Println("Usage: myFind -f -d -sl -ext extension path")
                os.Exit(1)
            }
            fileExt = os.Args[len(os.Args)-2]
        }
    }

	
	

    err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
		
        if info.IsDir() && findDirs {
            fmt.Println(path)
        }

        if info.Mode()&os.ModeSymlink != 0 && findSymlinks {
            linkPath, err := os.Readlink(path)
            if err != nil {
                fmt.Printf("%s -> [broken]\n", path)
            } else {
                fmt.Printf("%s -> %s\n", path, linkPath)
            }
        }

        if !info.IsDir() && !info.Mode().IsRegular() {
            return nil
        }

        if info.Mode().IsRegular() && findFiles {
            if fileExt == "" || strings.HasSuffix(info.Name(), "." + fileExt) {
                fmt.Println(path)
            }
        }

        return nil
    })

    if err != nil {
        fmt.Printf("Error occurred: %v\n", err)
        os.Exit(1)
    }
}
