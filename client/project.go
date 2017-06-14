package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/webdevwilson/tfwatch/model"
)

type Projects struct {
	sockAddr string
}

// NewProjectClient is used to create a project client
func NewProjectClient(sockAddr string) *Projects {
	return &Projects{sockAddr}
}

// Create creates a new project
func (p *Projects) Create(prj *model.Project) error {
	url := fmt.Sprintf("%s/api/projects", p.sockAddr)
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

// Update updates the project by guid
func (p *Projects) Update(prj *model.Project) error {
	url := fmt.Sprintf("%s/api/projects/%s", p.sockAddr, prj.GUID)

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

// Delete deletes a project by it's guid
func (p *Projects) Delete(guid string) error {
	url := fmt.Sprintf("%s/api/projects/%s", p.sockAddr, guid)

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

// Get fetches a project by it's guid
func (p *Projects) Get(guid string) (*model.Project, error) {
	url := fmt.Sprintf("%s/api/projects/%s", p.sockAddr, guid)

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

// List returns a list of projects
func (p *Projects) List() ([]model.Project, error) {
	url := fmt.Sprintf("%s/api/projects", p.sockAddr)

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
