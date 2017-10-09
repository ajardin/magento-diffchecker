package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	Patch   string
	Project string
)

func init() {
	flag.StringVar(&Patch, "patch", "", "Diff file to be applied on project.")
	flag.StringVar(&Project, "source", "", "Directory where project source files can be found.")
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, fileError := os.Open(Patch)
	CheckError(fileError)
	defer file.Close()

	parsePatchFile(file)
	//TODO
}

// parsePatchFile parses the .diff file and extracts modified files.
func parsePatchFile(file io.Reader) {
	pattern := regexp.MustCompile(`^Index: (?P<Path>\w+/(\w+/)*\w+\.(php|phtml|js))$`)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		matches := pattern.FindStringSubmatch(scanner.Text())
		if matches != nil {
			filePath := matches[1]

			if strings.HasSuffix(filePath, ".php") {
				//TODO
			} else if strings.HasSuffix(filePath, ".phtml") {
				//TODO
			} else if strings.HasSuffix(filePath, ".js") {
				//TODO
			}
		}
	}
}

// CheckError causes the current program to exit if an error occurred.
func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
