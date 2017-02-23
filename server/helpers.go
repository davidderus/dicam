package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type smallCamera struct {
	ID     string
	Port   string
	Active bool
}

type LayoutData struct {
	CamerasListing map[string]*smallCamera
}

type templateData struct {
	ViewData interface{}
	LayoutData
}

func getCameras() map[string]*smallCamera {
	camsRawList, _ := askClient("CAM-LIST")

	camerasList := map[string]*smallCamera{}

	for _, infosLine := range strings.Split(camsRawList, "\n") {
		infos := strings.Split(infosLine, " - ")

		key := strings.Replace(infos[0], "Cam. ", "", 1)

		isActive := infos[1] != "Not running"

		port := ""
		if isActive {
			port = strings.Replace(infos[2], "Port ", "", 1)
		}

		camerasList[key] = &smallCamera{ID: key, Port: port, Active: isActive}
	}

	return camerasList
}

func writeWithTemplate(response http.ResponseWriter, templateName string, templatePath string, data interface{}) {
	getTemplateDirForFile := func(file string) string {
		return filepath.Join("server", "templates", file)
	}

	templateFile, parseError := template.ParseFiles(
		getTemplateDirForFile("layout.html"),
		getTemplateDirForFile("navbar.html"),
		getTemplateDirForFile(templatePath),
	)

	if parseError != nil {
		log.Fatalf("Can't parse template for %s", templateName)
	}

	templateFile.ExecuteTemplate(response, "layout", templateData{data, LayoutData{CamerasListing: getCameras()}})
}
