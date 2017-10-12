package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type TemplatesDetails struct {
	pathList  []string
	themeList []string
}

var templates TemplatesDetails

// AnalyzeMagentoTemplates executes the Magento templates analysis.
func AnalyzeMagentoTemplates() {
	loadThemeList()

	for _, templatePath := range templates.pathList {
		isFromBaseDefault := strings.HasPrefix(templatePath, filepath.FromSlash("app/design/frontend/base/default/"))
		isFromEnterpriseDefault := strings.HasPrefix(templatePath, filepath.FromSlash("app/design/frontend/enterprise/default/"))

		if !isFromBaseDefault && !isFromEnterpriseDefault {
			continue
		}

		for _, themeName := range templates.themeList {
			checkTemplateCompleteOverride(templatePath, themeName)
		}
	}
}

// loadThemeList loads the theme list from the project directory.
func loadThemeList() {
	packageEntries, packageReadError := ioutil.ReadDir(Project + filepath.FromSlash("/app/design/frontend"))
	CheckError(packageReadError)

	for _, packageEntry := range packageEntries {
		if !packageEntry.IsDir() || packageEntry.Name() == "base" || packageEntry.Name() == "enterprise" {
			continue
		}

		themeEntries, themeReadError := ioutil.ReadDir(Project + filepath.FromSlash("/app/design/frontend/") + packageEntry.Name())
		CheckError(themeReadError)

		for _, themeEntry := range themeEntries {
			if !themeEntry.IsDir() {
				continue
			}
			templates.themeList = append(templates.themeList, packageEntry.Name()+string(os.PathSeparator)+themeEntry.Name())
		}
	}
}

// checkTemplateCompleteOverride checks whether a complete override has been performed.
func checkTemplateCompleteOverride(templatePath string, themeName string) {
	targetPath := ""

	if strings.Contains(templatePath, filepath.FromSlash("enterprise/default")) {
		targetPath = strings.Replace(templatePath, filepath.FromSlash("enterprise/default"), themeName, -1)
	} else {
		targetPath = strings.Replace(templatePath, filepath.FromSlash("base/default"), themeName, -1)
	}

	if IsFileExists(Project + string(os.PathSeparator) + targetPath) {
		color.Set(color.FgRed)
		fmt.Println("Complete override detected (template):", templatePath, "==>", targetPath)
		color.Unset()
	}
}
