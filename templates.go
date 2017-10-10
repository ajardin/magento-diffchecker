package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
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
		if !strings.HasPrefix(templatePath, "app/design/frontend/base/default/") && !strings.HasPrefix(templatePath, "app/design/frontend/enterprise/default/") {
			continue
		}

		for _, themeName := range templates.themeList {
			checkTemplateCompleteOverride(templatePath, themeName)
		}
	}
}

// loadThemeList loads the theme list from the project directory.
func loadThemeList() {
	packageEntries, packageReadError := ioutil.ReadDir(Project + "/app/design/frontend")
	CheckError(packageReadError)

	for _, packageEntry := range packageEntries {
		if !packageEntry.IsDir() || packageEntry.Name() == "base" || packageEntry.Name() == "enterprise" {
			continue
		}

		themeEntries, themeReadError := ioutil.ReadDir(Project + "/app/design/frontend/" + packageEntry.Name())
		CheckError(themeReadError)

		for _, themeEntry := range themeEntries {
			if !themeEntry.IsDir() {
				continue
			}
			templates.themeList = append(templates.themeList, packageEntry.Name()+"/"+themeEntry.Name())
		}
	}
}

// checkTemplateCompleteOverride checks whether a complete override has been performed.
func checkTemplateCompleteOverride(templatePath string, themeName string) {
	targetPath := ""

	if strings.Contains(templatePath, "enterprise/default") {
		targetPath = strings.Replace(templatePath, "enterprise/default", themeName, -1)
	} else {
		targetPath = strings.Replace(templatePath, "base/default", themeName, -1)
	}

	if IsFileExists(Project + "/" + targetPath) {
		color.Set(color.FgRed)
		fmt.Println("Complete override detected (template):", templatePath, "==>", targetPath)
		color.Unset()
	}
}
