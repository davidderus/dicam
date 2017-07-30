package server

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

type smallCamera struct {
	ID     string
	Port   string
	Active bool
}

// LayoutData contains all data required for the persisted layout display
type LayoutData struct {
	CamerasListing map[string]*smallCamera
}

type templateBindedData struct {
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

func writeWithTemplate(response http.ResponseWriter, templateName string, templatePath string, data interface{}) error {
	templates := template.New("")

	var (
		templateDir  string
		templateData []byte
		assetError   error
	)

	templatesToRender := []string{"layout.html", "navbar.html", templatePath}

	for _, templateName := range templatesToRender {
		templateDir = filepath.Join("server", "templates", templateName)
		templateData, assetError = Asset(templateDir)

		// Skipping template on asset error
		if assetError == nil {
			templates.New(templateDir).Parse(string(templateData))
		}
	}

	templates.ExecuteTemplate(response, "layout", templateBindedData{data, LayoutData{CamerasListing: getCameras()}})

	return nil
}
