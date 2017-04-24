package main

import (
	"log"
	"path"
	"path/filepath"

	"github.com/webdevwilson/terraform-ci/config"
	"github.com/webdevwilson/terraform-ci/model"
	"github.com/webdevwilson/terraform-ci/routes"
	_ "github.com/webdevwilson/terraform-ci/task"
)

func main() {
	settings := config.Get()
	port := settings.Port

	// this should be temporary, until we can get a proper CRUD frontend
	bootstrapProjects(settings.CheckoutDirectory)

	go routes.StartServer(port)

	for {
	}
}

func bootstrapProjects(dir string) {
	dirs, err := filepath.Glob(path.Join(dir, "*"))
	if err != nil {
		log.Fatalf("[FATAL] Error encountered loading projects: %s", err)
	}

	for _, dir := range dirs {
		name := path.Base(dir)
		prj, err := model.GetProjectByName(name)

		if err != nil {
			log.Fatalf("[FATAL] Error bootstrapping projects: %s", err)
		}

		if prj == nil {
			model.CreateProject(&model.Project{
				Name: name,
			})
		}
	}
}
