/*******************************************************************************
 * Copyright (c) 2019 IBM Corporation and others.
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v20.html
 *
 * Contributors:
 *     IBM Corporation - initial API and implementation
 *******************************************************************************/

package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Template represents a template.
type Template struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	Language    string `json:"language"`
	URL         string `json:"url"`
	ProjectType string `json:"projectType"`
}

// TemplateDescription represents a description of a template.
type TemplateDescription struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Location    string `json:"location"`
	ProjectType string `json:"projectType"`
}

// TemplateRepo represents a template repository.
type TemplateRepo struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

// ListTemplatesCommand lists all templates Codewind is aware of.
func ListTemplatesCommand() {
	// TODO: find this dynamically
	const pathToTemplateRepoFile = "../codewind-installer/resources/repository_list.json"
	templates, err := GetTemplates(pathToTemplateRepoFile)
	if err != nil {
		fmt.Printf("Error reading templates at %s: %q", pathToTemplateRepoFile, err)
		return
	}
	fmt.Println("Templates available:")
	printTemplates(templates)
}

// GetTemplates extracts URLs from the file at the provided path,
// then gets the template descriptions from those URLs,
// then formats the descriptions into template objects.
func GetTemplates(pathToTemplateRepoFile string) ([]Template, error) {
	var templates []Template

	templateRepos, err := ReadJSONFile(pathToTemplateRepoFile)
	if err != nil {
		return nil, err
	}

	for _, repoDetails := range templateRepos {
		templateDescriptions, _ := GetTemplateDescriptions(repoDetails)
		if err != nil {
			return nil, err
		}

		for _, templateDescription := range templateDescriptions {
			template := FormatTemplate(templateDescription)
			templates = append(templates, template)
		}
	}

	return templates, nil
}

// ReadJSONFile reads the JSON file at the provided filepath.
func ReadJSONFile(filepath string) ([]TemplateRepo, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteArray, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var repoList []TemplateRepo
	json.Unmarshal(byteArray, &repoList)

	return repoList, nil
}

// GetTemplateDescriptions gets TemplateDescriptions by making GET requests to urls found in the provided repoDetails.
func GetTemplateDescriptions(repoDetails TemplateRepo) ([]TemplateDescription, error) {
	resp, err := http.Get(repoDetails.URL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var templateDescriptions []TemplateDescription
	json.Unmarshal(byteArray, &templateDescriptions)

	return templateDescriptions, nil
}

// FormatTemplate extracts fields from a TemplateDescription and inserts them into a Template.
func FormatTemplate(templateDesc TemplateDescription) Template {
	// TODO: try func (templateDesc *TemplateDescription) Format() Template {
	template := Template{
		Label:       templateDesc.DisplayName,
		Description: templateDesc.Description,
		Language:    templateDesc.Language,
		URL:         templateDesc.Location,
		ProjectType: templateDesc.ProjectType,
	}
	return template
}

// TODO: try .Print()
func printTemplates(templates []Template) {
	for _, template := range templates {
		PrettyPrintJSON(template)
	}
}

// PrettyPrintJSON prints JSON prettily.
func PrettyPrintJSON(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(s))
}
