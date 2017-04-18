package routes

import (
	"fmt"
	"os"
	"testing"

	"net/http"

	"encoding/json"

	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/webdevwilson/terraform-ui/config"
	"github.com/webdevwilson/terraform-ui/model"
)

var prjs = []model.Project{
	model.Project{
		"",
		"foo",
		"git@github.com:webdevwilson/terraform-ci",
		"/foo",
	},
}

var urlBase string

func TestMain(m *testing.M) {
	go StartServer(config.Get().Port)
	urlBase = fmt.Sprintf("http://localhost:%s", config.Get().Port)
	result := m.Run()
	os.Exit(result)
}

func Test_Create(t *testing.T) {
	url := fmt.Sprintf("%s/api/projects", urlBase)
	body, err := json.Marshal(prjs[0])
	if err != nil {
		t.Error(err)
	}

	// make request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}

	resp, err := http.DefaultClient.Do(req)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "", resp.Body)
}
