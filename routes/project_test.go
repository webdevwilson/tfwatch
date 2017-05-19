package routes

import (
	"fmt"
	"math/rand"
	"testing"

	"net/http"

	"encoding/json"

	"bytes"

	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/webdevwilson/terraform-ci/model"
)

var prjs = []model.Project{
	model.Project{
		"",
		"foo",
		"git@github.com:webdevwilson/terraform-ci",
		"/foo",
		map[string]string{
			"FOO": "BAR",
		},
	},
	model.Project{
		"",
		"bar",
		"git@github.com:webdevwilson/terraform-ci",
		"/bar",
		map[string]string{},
	},
}

func StartTestServer() string {
	// get random ephemeral port
	port := rand.Intn(16383) + 49152
	go StartServer(port)
	return fmt.Sprintf("http://localhost:%d", port)
}

func Create(sockAddr string, prj *model.Project) error {
	url := fmt.Sprintf("%s/api/projects", sockAddr)
	body, err := json.Marshal(prj)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Invalid status code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(prj)
	if err != nil {
		return err
	}

	return nil
}

func Update(sockAddr string, prj *model.Project) error {
	url := fmt.Sprintf("%s/api/projects/%s", sockAddr, prj.GUID)

	body, err := json.Marshal(prj)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Invalid status code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(prj)
	if err != nil {
		return err
	}

	return nil
}

func Delete(sockAddr string, guid string) error {
	url := fmt.Sprintf("%s/api/projects/%s", sockAddr, guid)

	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(""))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Invalid status code %d", resp.StatusCode)
	}

	return nil
}

func Get(sockAddr string, guid string) (*model.Project, error) {
	url := fmt.Sprintf("%s/api/projects/%s", sockAddr, guid)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid status code %d", resp.StatusCode)
	}

	var result = &model.Project{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func List(sockAddr string) ([]model.Project, error) {
	url := fmt.Sprintf("%s/api/projects", sockAddr)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid status code %d", resp.StatusCode)
	}

	var result []model.Project
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Test_Project_API(t *testing.T) {
	sockAddr := StartTestServer()

	for i, v := range prjs {
		err := Create(sockAddr, &v)
		prjs[i].GUID = v.GUID
		if err != nil {
			t.Error(err)
		}
	}

	prj, err := Get(sockAddr, prjs[0].GUID)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.NotEmpty(t, prj.GUID)
	assert.Equal(t, prjs[0].Name, prj.Name)
	assert.Equal(t, prjs[0].RepoURL, prj.RepoURL)
	assert.Equal(t, prjs[0].RepoPath, prj.RepoPath)
	assert.Equal(t, prjs[0].Settings["FOO"], prj.Settings["FOO"])

	prjs[1].Name = "abc"
	guid := prjs[1].GUID
	err = Update(sockAddr, &prjs[1])
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, guid, prjs[1].GUID)
	assert.Equal(t, "abc", prjs[1].Name)

	// do list
	list, err := List(sockAddr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, len(prjs), len(list))

	// do delete
	err = Delete(sockAddr, prjs[1].GUID)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// verify delete
	list, err = List(sockAddr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	assert.Equal(t, len(prjs)-1, len(list))

}
