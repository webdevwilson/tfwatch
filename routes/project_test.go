package routes

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webdevwilson/terraform-ci/client"
	"github.com/webdevwilson/terraform-ci/controller"
	"github.com/webdevwilson/terraform-ci/execute"
	"github.com/webdevwilson/terraform-ci/model"
	"github.com/webdevwilson/terraform-ci/persist"
	"github.com/webdevwilson/terraform-ci/test"
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
	stateDir := path.Join(checkoutDir, ".terraform-ci")
	logDir := path.Join(stateDir, "logs")

	store, _ := persist.NewBoltStore(stateDir)
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

func Test_Project_API(t *testing.T) {
	sockAddr := startTestServer()
	projects := client.NewProjectClient(sockAddr)

	for i, v := range prjs {
		err := projects.Create(&v)
		prjs[i].GUID = v.GUID
		if err != nil {
			t.Error(err)
		}
	}

	prj, err := projects.Get(prjs[0].GUID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.NotEmpty(t, prj.GUID)
	assert.Equal(t, prjs[0].Name, prj.Name)
	assert.Equal(t, prjs[0].Settings["FOO"], prj.Settings["FOO"])

	prjs[1].Name = "abc"
	guid := prjs[1].GUID
	err = projects.Update(&prjs[1])
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, guid, prjs[1].GUID)
	assert.Equal(t, "abc", prjs[1].Name)

	// do list
	list, err := projects.List()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, len(prjs)+3, len(list))

	// do delete
	err = projects.Delete(prjs[1].GUID)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// verify delete
	list, err = projects.List()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, len(prjs)+2, len(list))

}
