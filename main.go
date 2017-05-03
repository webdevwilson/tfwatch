package main

import (
	"log"
	"path"
	"path/filepath"

	"github.com/webdevwilson/terraform-ci/config"
	_ "github.com/webdevwilson/terraform-ci/execute"
	"github.com/webdevwilson/terraform-ci/model"
	"github.com/webdevwilson/terraform-ci/routes"
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
	log.Printf("[INFO] Bootstrapping projects in %s", dir)

	dirs, err := filepath.Glob(path.Join(dir, "*"))
	if err != nil {
		log.Fatalf("[FATAL] Error encountered loading projects: %s", err)
		return
	}

	for _, dir := range dirs {

		// does this path have .tf files?
		tf, err := filepath.Glob(path.Join(dir, "*.tf"))
		if err != nil {
			log.Fatalf("[FATAL] Error encountered searching for .tf files: %s", err)
			return
		}

		// none found? continue to next
		if len(tf) == 0 {
			log.Printf("[INFO] No .tf files found in %s, skipping", dir)
			continue
		}

		name := path.Base(dir)
		prj, err := model.GetProjectByName(name)

		if err != nil {
			log.Fatalf("[FATAL] Error bootstrapping projects: %s", err)
		}

		if prj == nil {
			log.Printf("[INFO] Creating project '%s' found at '%s'", name, dir)
			model.CreateProject(&model.Project{
				Name:      name,
				LocalPath: dir,
			})
		}
	}
}
