package model

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webdevwilson/tfwatch/persist"
	"github.com/webdevwilson/tfwatch/test"
)

var data = []Project{
	Project{
		GUID:      "0000-0000",
		Name:      "terraform_noplan",
		LocalPath: "../fixtures/terraform_noplan",
		Settings:  map[string]string{}},
	Project{
		GUID:      "1111-1111",
		Name:      "terraform_planned",
		LocalPath: "../fixtures/terraform_planned",
		Settings:  map[string]string{}},
}

func TestMain(m *testing.M) {
	test.SuppressLogs()
	os.Exit(m.Run())
}

func TestProject_get(t *testing.T) {
	lfs, err := persist.NewLocalFileStore("./")
	assert.NoError(t, err)

	var prj Project
	lfs.Get("projects", "76e29370-f046-491b-5b01-4bbd4e17aabc", &prj)
}

func TestPlan_no_plan(t *testing.T) {
	_, err := data[0].Plan()

	assert.Error(t, err, "")
}

func TestPlan_plan(t *testing.T) {
	plan, err := data[1].Plan()
	assert.Nil(t, err)
	assert.NotNil(t, plan)
}

func TestPlan_changes(t *testing.T) {
	plan, err := data[1].Plan()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(plan.ResourceChanges()))

	change := plan.ResourceChanges()[0]

	assert.Equal(t, "root.local_file.file", change.ResourceID)
	assert.Equal(t, "Create", change.Action)
}
