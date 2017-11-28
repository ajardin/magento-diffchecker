package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
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
	flag.StringVar(&Patch, "patch", os.Getenv("PATCH"), "Diff file to be applied on project.")
	flag.StringVar(&Project, "project", os.Getenv("PROJECT"), "Directory where project source files can be found.")
}

func main() {
	loadFlagValues()

	file, fileError := os.Open(Patch)
	CheckError(fileError)
	defer file.Close()

	color.Set(color.FgGreen)
	fmt.Println("Magento analysis in progress...")
	color.Unset()

	fmt.Println("PATCH =", Patch)
	fmt.Println("PROJECT =", Project)

	parsePatchFile(file)
	AnalyzeMagentoClasses()
	AnalyzeMagentoTemplates()

	color.Set(color.FgGreen)
	fmt.Println("Magento analysis successfully completed.")
	color.Unset()
}

// loadFlagValues loads and checks whether all flag values are valid.
func loadFlagValues() {
	flag.Parse()

	if IsFileExists(Patch) != true {
		log.Fatal("The 'patch' flag must be a valid file.")
	}

	if IsDirectoryExists(Project) != true {
		log.Fatal("The 'project' flag must be a valid directory.")
	}
}

// parsePatchFile parses the .diff file and extracts modified files.
func parsePatchFile(file io.Reader) {
	pattern := regexp.MustCompile(`^diff --git (a/)?(?P<Path>\w+/(\w+/)*\w+\.(php|phtml))`)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		matches := pattern.FindStringSubmatch(scanner.Text())
		if matches != nil {
			filePath := matches[2]

			if strings.HasSuffix(filePath, ".php") {
				classes.pathList = append(classes.pathList, filePath)
			} else if strings.HasSuffix(filePath, ".phtml") {
				templates.pathList = append(templates.pathList, filePath)
			}
		}
	}

	sort.Strings(classes.pathList)
	sort.Strings(templates.pathList)
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
