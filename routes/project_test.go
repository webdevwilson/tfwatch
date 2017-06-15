package routes

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"github.com/webdevwilson/tfwatch/client"
	"github.com/webdevwilson/tfwatch/controller"
	"github.com/webdevwilson/tfwatch/execute"
	"github.com/webdevwilson/tfwatch/model"
	"github.com/webdevwilson/tfwatch/persist"
	"github.com/webdevwilson/tfwatch/test"
)

var project = model.Project{
	GUID:      "001",
	Name:      "foo",
	LocalPath: "/foo",
	Settings: map[string]string{
		"FOO": "BAR",
	},
	Status: model.ProjectStatusNew,
}

// Ideally, we should just build a context, however there are cyclic dependency issues
// so I am sticking this here for now
func startTestServer() string {

	// get random ephemeral port
	test.SuppressLogs()
	port := uint16(rand.Int31n(16383) + 49152)

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// use fixtures directory
	checkoutDir := path.Clean(path.Join(cwd, "..", "fixtures"))
	siteDir := path.Clean(path.Join(cwd, "..", "site", "dist"))
	stateDir := path.Join(checkoutDir, ".tfwatch")
	logDir := path.Join(stateDir, "logs")

	store, err := persist.NewBoltStore(stateDir)
	if err != nil {
		panic(err)
	}
	exec := execute.NewExecutor(store, logDir)
	sys := controller.NewSystemController([]*controller.SystemConfigurationValue{}, exec)
	prj := controller.NewProjectsController(checkoutDir, store, exec, 5*time.Minute, false)

	server := InitializeServer(port, ioutil.Discard, sys, prj, siteDir)
	go server.Start()

	return fmt.Sprintf("http://localhost:%d", port)
}

func Test_Project_Create(t *testing.T) {
	sockAddr := startTestServer()
	projects := client.NewProjectClient(sockAddr)

	err := projects.Create(&project)
	if err != nil {
		t.Error(err)
	}

	projects.Delete(project.GUID)
}
