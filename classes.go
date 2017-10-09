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

var (
	ClassePathList []string
	rewriteList    map[string]string
)

// AnalyzeMagentoClasses executes the Magento classes analysis.
func AnalyzeMagentoClasses() {
	rewriteList = make(map[string]string)
	filepath.Walk(Project+"/app/code/local/", loadRewrites())

	for _, classPath := range ClassePathList {
		if classPath == "app/Mage.php" || classPath == "app/code/core/Mage/Core/functions.php" {
			continue
		}

		absolutePath := Project + "/" + classPath
		checkCompleteOverride(absolutePath)

		if IsFileExists(absolutePath) {
			checkPartialOverride(absolutePath)
		}
	}
}

// loadRewrites loads all rewrite rules from config.xml files in local pool.
func loadRewrites() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.Name() == "config.xml" {
			extractRewritesFromFile(path)
		}
		return nil
	}
}

// extractRewritesFromFile extracts rewrites rules described in the given config.xml file.
func extractRewritesFromFile(path string) {
	file, openError := os.Open(path)
	CheckError(openError)
	defer file.Close()

	document, parseError := xmlquery.Parse(file)
	CheckError(parseError)

	for _, node := range xmlquery.Find(document, "//rewrite/*") {
		className := buildFullClassName(node.Parent.Parent.Data, node.Parent.Parent.Parent.Data, node.Data)
		rewriteList[className] = node.InnerText()
	}
}

// buildFullClassName builds the full Magento class (Packagename_Modulename_Classtype_Classname).
func buildFullClassName(moduleName string, classType string, className string) string {
	tmp := ""

	//TODO: check community modules /!\
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

// checkCompleteOverride checks whether a complete override has been performed.
func checkCompleteOverride(sourcePath string) {
	replacer := strings.NewReplacer("core", "local", "community", "local")
	targetPath := replacer.Replace(sourcePath)

	if IsFileExists(sourcePath) && IsFileExists(targetPath) {
		originalRelativePath := strings.TrimPrefix(sourcePath, Project+"/")
		customRelativePath := strings.TrimPrefix(targetPath, Project+"/")

		color.Set(color.FgRed)
		fmt.Println("Complete override detected:", originalRelativePath, "==>", customRelativePath)
		color.Unset()
	}
}

// checkPartialOverride checks whether a partial override has been performed.
func checkPartialOverride(sourcePath string) {
	className := extractClassNameFromFile(sourcePath)
	if className != "" {
		value, isset := rewriteList[className]
		if isset == true {
			color.Set(color.FgYellow)
			fmt.Println("Partial override detected:", className, "==>", value)
			color.Unset()
		}
	}
}

// extractClassNameFromFile extracts the PHP class name from the given source path.
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
