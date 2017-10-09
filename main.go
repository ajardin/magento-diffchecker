package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	Patch   string
	Project string
)

func init() {
	flag.StringVar(&Patch, "patch", "", "Diff file to be applied on project.")
	flag.StringVar(&Project, "project", "", "Directory where project source files can be found.")
}

func main() {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	checkFlags()

	file, fileError := os.Open(Patch)
	CheckError(fileError)
	defer file.Close()

	parsePatchFile(file)
	AnalyzeMagentoClasses()
	//TODO
}

// checkFlags checks whether all flags are valid.
func checkFlags() {
	if IsFileExists(Patch) != true {
		log.Fatal("The 'patch' flag must be a valid file.")
	}

	if IsDirectoryExists(Project) != true {
		log.Fatal("The 'project' flag must be a valid directory.")
	}
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
				ClassePathList = append(ClassePathList, filePath)
			} else if strings.HasSuffix(filePath, ".phtml") {
				//TODO
			} else if strings.HasSuffix(filePath, ".js") {
				//TODO
			}
		}
	}

	sort.Strings(ClassePathList)
}

// CheckError causes the current program to exit if an error occurred.
func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// IsFileExists checks whether the file exists.
func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectoryExists checks whether the directory exists.
func IsDirectoryExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}
