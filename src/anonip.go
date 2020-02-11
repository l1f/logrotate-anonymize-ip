// This software anonymizes log files using the parameter specified by the flags
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
)

var (
	logFilePath  string
	tmpFilePath  string
	IPBlock      = "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	regexPattern = IPBlock + "\\." + IPBlock + "\\." + IPBlock + "\\." + IPBlock
)

func init() {
	handleFlags()
}

// extracts the arguments passed via the flags and processes these
func handleFlags()  {
	help := flag.Bool("h", false, "Show help message")
	flag.StringVar(&logFilePath, "f", "", "Which file should be anonymized?")
	flag.StringVar(&tmpFilePath, "t", "", "Where to save the temp file?")
	flag.Parse()

	if *help {
		printHelp()
	}

	if logFilePath == "" {
		printHelp()
	}

	// If the temp file arguments is not specified, the pid of the current process is used as default.
	if tmpFilePath == "" {
		tmpFilePath = "/tmp/anon_" + strconv.Itoa(os.Getpid()) + "_.tmp"
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

// Opens the two files, reads the original file and then overwrites it with the anonymous logs.
func main() {
	logFile, err := os.Open(logFilePath)
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

	logFile.Close() // To be on the safe side, the file is closed before the function ends, so that nothing goes wrong when overwriting.

	err = copyFile(tmpFilePath, logFilePath)
	if err != nil {
		log.Fatal(err)
	}
}

// Accepts the current line and replaces the IPs in the string.
func replaceIP(input string) string {
	regEx := regexp.MustCompile(regexPattern)
	return regEx.ReplaceAllString(input, "[ Anonymized IP ]")
}

// Overwrite the src file via the destionation
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}