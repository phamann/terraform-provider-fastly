//go:generate go run parse-templates.go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

type Page struct {
	name         string
	path         string
	Data         PageData
	DataMenu     []MenuItem
	ResourceMenu []MenuItem
}

type PageData struct {
	ServiceType string
}

type MenuItem struct {
	Lookup string
	Link   string
	Title  string
}

func main() {
	baseDir := getBaseDir()
	tmplDir := baseDir + "/docs_src/"
	docsDir := baseDir + "/docs/"

	var dataPages = []Page{
		{
			name: "ip_ranges",
			path: docsDir + "data-sources/ip_ranges.html.markdown",
		},
		{
			name: "waf_rules",
			path: docsDir + "data-sources/waf_rules.html.markdown",
		},
	}

	var resourcePages = []Page{
		{
			name: "service_v1",
			path: docsDir + "resources/service_v1.html.markdown",
			Data: PageData{
				"vcl",
			},
		},
		{
			name: "service_compute",
			path: docsDir + "resources/service_compute.html.markdown",
			Data: PageData{
				"wasm",
			},
		},
		{
			name: "service_dictionary_items_v1",
			path: docsDir + "resources/service_dictionary_items_v1.html.markdown",
		},
		{
			name: "service_acl_entries_v1",
			path: docsDir + "resources/service_acl_entries_v1.html.markdown",
		},
		{
			name: "service_dynamic_snippet_content_v1",
			path: docsDir + "resources/service_dynamic_snippet_content_v1.html.markdown",
		},
		{
			name: "service_waf_configuration",
			path: docsDir + "resources/service_waf_configuration.html.markdown",
		},
		{
			name: "user_v1",
			path: docsDir + "resources/user_v1.html.markdown",
		},
	}

	var pages = append(resourcePages, Page{
		name:         "fastly_erb",
		path:         docsDir + "fastly.erb",
		DataMenu:     generateMenuItems("data-sources", dataPages),
		ResourceMenu: generateMenuItems("resources", resourcePages),
	})

	renderPages(getTemplate(tmplDir), pages)
}

func generateMenuItems(pageType string, pages []Page) []MenuItem {
	var menuItems []MenuItem
	for _, p := range pages {
		menuItems = append(menuItems, MenuItem{
			Lookup: fmt.Sprintf("docs-fastly-resource-%s", strings.ReplaceAll(p.name, "_", "-")),
			Link:   fmt.Sprintf("/providers/fastly/%s/%s.html", pageType, p.name),
			Title:  fmt.Sprintf("fastly_%s", p.name),
		})
	}
	return menuItems
}

func renderPages(t *template.Template, pages []Page) {
	for _, p := range pages {
		renderPage(t, p)
	}
}

func renderPage(t *template.Template, p Page) {
	f, err := os.Create(p.path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = t.ExecuteTemplate(f, p.name, p)
	if err != nil {
		panic(err)
	}
}

func getTemplate(tmplDir string) *template.Template {
	var templateFiles []string
	filepath.Walk(tmplDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".tmpl" {
			templateFiles = append(templateFiles, path)
		}
		return nil
	})
	return template.Must(template.ParseFiles(templateFiles...))
}

func getBaseDir() string {
	_, scriptPath, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Could not get current working directory")
	}
	tpgDir := scriptPath
	for !strings.HasPrefix(filepath.Base(tpgDir), "terraform-provider-") && tpgDir != "/" {
		tpgDir = filepath.Clean(tpgDir + "/..")
	}
	if tpgDir == "/" {
		log.Fatal("Script was run outside of fastly provider directory")
	}
	return tpgDir
}
