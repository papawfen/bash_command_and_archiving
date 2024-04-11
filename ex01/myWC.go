package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Flags struct {
	lineCount bool
	charCount bool
	wordCount bool
}

func printError(err string) {
	fmt.Println(err)
	os.Exit(1)
}

func countChars(path string, file *os.File) string {
	var charCount int64
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		charCount++
	}
	return strconv.FormatInt(charCount, 10) + " " + path
}

func countlines(path string, file *os.File) string {
	var lineCount int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount += 1
	}
	return strconv.FormatInt(lineCount - 1, 10) + " " + path
}

func countWords(path string, file *os.File) string {
	var wordCount int64
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordCount++
	}
	return strconv.FormatInt(wordCount, 10) + " " + path
}

func countElements(path string, flags Flags) string {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		printError("cant open file")
	}
	defer file.Close()
	var res string 
	if flags.charCount {
		res = countChars(path, file)
	} else if flags.lineCount {
		res = countChars(path, file)
	} else if flags.wordCount {
		res = countWords(path, file)
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myWC [-l] || [-m] || [-w] path")
		os.Exit(1)
	}
	var flags Flags
	var noFlags bool
	noFlags = false
	switch os.Args[1] {
		case "-l":
			flags.lineCount = true
		case "-m":
			flags.charCount = true
		case "-w":
			flags.wordCount = true
		default:
			noFlags = true
			flags.wordCount = true
	}
	res := make(chan string)
	var startCount int
	if noFlags {
		startCount = 1
	} else {
		startCount = 2
	}
	for _, arg := range os.Args[startCount:] {
		go func(arg string) {
			line := countElements(arg, flags)
			res <- line
		}(arg)
	}
	for range os.Args[startCount:] {
		fmt.Println("\t", <-res)
	}
}