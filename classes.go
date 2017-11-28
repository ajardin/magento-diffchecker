package main

import (
	"bufio"
	"fmt"
	"github.com/antchfx/xquery/xml"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ClassesDetails struct {
	pathList    []string
	rewriteList map[string]string
}

var classes ClassesDetails

// AnalyzeMagentoClasses executes the Magento classes analysis.
func AnalyzeMagentoClasses() {
	classes.rewriteList = make(map[string]string)
	filepath.Walk(Project+filepath.FromSlash("/app/code/local/"), loadClassRewrites())

	for _, classPath := range classes.pathList {
		isMageFile := classPath == filepath.FromSlash("app/Mage.php")
		isFunctionsFile := classPath == filepath.FromSlash("app/code/core/Mage/Core/functions.php")

		isDownloaderFile := strings.HasPrefix(classPath, "downloader"+string(os.PathSeparator))
		isErrorFile := strings.HasPrefix(classPath, "errors"+string(os.PathSeparator))
		isLibFile := strings.HasPrefix(classPath, "lib"+string(os.PathSeparator))

		if isMageFile || isFunctionsFile || isDownloaderFile || isErrorFile || isLibFile {
			continue
		}

		absolutePath := Project + string(os.PathSeparator) + classPath
		checkClassCompleteOverride(absolutePath)

		if IsFileExists(absolutePath) {
			checkClassPartialOverride(absolutePath)
		}
	}
}

// loadClassRewrites loads all rewrite rules from config.xml files in local pool.
func loadClassRewrites() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.Name() == "config.xml" {
			extractClassRewritesFromFile(path)
		}
		return nil
	}
}

// extractClassRewritesFromFile extracts rewrites rules described in the config.xml file.
func extractClassRewritesFromFile(path string) {
	file, openError := os.Open(path)
	CheckError(openError)
	defer file.Close()

	document, parseError := xmlquery.Parse(file)
	CheckError(parseError)

	for _, node := range xmlquery.Find(document, "//rewrite/*") {
		className := buildFullClassName(node.Parent.Parent.Data, node.Parent.Parent.Parent.Data, node.Data)
		classes.rewriteList[className] = node.InnerText()
	}
}

// buildFullClassName builds the full Magento class (Packagename_Modulename_Classtype_Classname).
func buildFullClassName(moduleName string, classType string, className string) string {
	tmp := ""

	//TODO: check controllers & community modules /!\
	if strings.Contains(className, "enterprise") {
		tmp = moduleName + "_" + strings.TrimSuffix(classType, "s") + "_" + className
	} else {
		tmp = "mage" + "_" + moduleName + "_" + strings.TrimSuffix(classType, "s") + "_" + className
	}

	parts := strings.Split(tmp, "_")
	for idx, val := range parts {
		parts[idx] = strings.Title(val)
	}

	return strings.Join(parts, "_")
}

// checkClassCompleteOverride checks whether a complete override has been performed.
func checkClassCompleteOverride(sourcePath string) {
	replacer := strings.NewReplacer("core", "local", "community", "local")
	targetPath := replacer.Replace(sourcePath)

	if IsFileExists(sourcePath) && IsFileExists(targetPath) {
		originalRelativePath := strings.TrimPrefix(sourcePath, Project+string(os.PathSeparator))
		customRelativePath := strings.TrimPrefix(targetPath, Project+string(os.PathSeparator))

		color.Set(color.FgRed)
		fmt.Println("Complete override detected (class):", originalRelativePath, "==>", customRelativePath)
		color.Unset()
	}
}

// checkClassPartialOverride checks whether a partial override has been performed.
func checkClassPartialOverride(sourcePath string) {
	className := extractClassNameFromFile(sourcePath)
	if className != "" {
		value, isset := classes.rewriteList[className]
		if isset == true {
			color.Set(color.FgYellow)
			fmt.Println("Partial override detected (class):", className, "==>", value)
			color.Unset()
		}
	}
}

// extractClassNameFromFile extracts the PHP class name from the source path.
func extractClassNameFromFile(sourcePath string) string {
	result := ""

	file, openError := os.Open(sourcePath)
	CheckError(openError)
	defer file.Close()

	pattern := regexp.MustCompile(`(class|interface)\s+(?P<ClassName>\w+)`)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		matches := pattern.FindStringSubmatch(scanner.Text())
		if matches != nil {
			result = matches[2]
			break
		}
	}

	return result
}
