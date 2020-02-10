package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)


var filePath string
var tmpFilePath string

func init() {
	handleFlags()
}

func handleFlags()  {
	help := flag.Bool("h", false, "Show help message")
	flag.StringVar(&filePath, "f", "", "Which file should be anonymized?")
	flag.StringVar(&tmpFilePath, "t", "", "Where to save the temp file?")
	flag.Parse()

	if *help {
		printHelp()
	}

	if filePath == "" {
		printHelp()
	}

	if tmpFilePath == "" {
		timeString := strconv.Itoa(int(time.Now().Unix()))
		tmpFilePath = "/tmp/anon_" + timeString + "_.tmp"
	}
}

func printHelp() {
	fmt.Println("Anon-IP")
	fmt.Println()
	fmt.Printf("Usage: %s [OPTION]\n", os.Args[0])
	fmt.Println()
	fmt.Println("Options: ")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	logFile, err := os.Open(filePath)
	if err != nil{
		log.Fatal(err)
	}

	tmpFile, err := os.Create(tmpFilePath)
	if err != nil{
		log.Fatal(err)
	}

	defer func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			log.Fatal(err)
		}
	}()

	writer := bufio.NewWriter(tmpFile)
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		_, err := writer.WriteString(replaceIP(scanner.Text() + "\n"))
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	logFile.Close()

	err = CopyFile(tmpFilePath, filePath)
	if err != nil {
		log.Fatal(err)
	}
}

func replaceIP(input string) string {
	numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock

	regEx := regexp.MustCompile(regexPattern)
	ip := regEx.FindString(input)
	if ip == "" {
		return input
	}

	return regEx.ReplaceAllString(input, "[ Anonymized IP ]")
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}