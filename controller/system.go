package controller

import (
	"log"

	"strings"

	"github.com/webdevwilson/tfwatch/execute"
)

type SystemConfigurationValue struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type System interface {
	GetConfiguration() []*SystemConfigurationValue
}

type systemController struct {
	config   []*SystemConfigurationValue
	executor execute.Executor
}

// NewSystemController creates a new controller for working with system information
func NewSystemController(config []*SystemConfigurationValue, executor execute.Executor) System {

	sys := &systemController{
		config:   config,
		executor: executor,
	}

	go sys.terraformVersion()

	return sys
}

func (s *systemController) GetConfiguration() []*SystemConfigurationValue {
	return s.config
}

func (s *systemController) terraformVersion() {
	st, err := s.executor.Schedule(&execute.Task{
		Command: "terraform",
		Args:    []string{"version"},
	})

	if err != nil {
		log.Printf("[ERROR] Error getting terraform version: %s", err)
	}

	r := <-st.Channel

	s.config = append(s.config, &SystemConfigurationValue{
		ID:    "TerraformVersion",
		Name:  "Terraform Version",
		Value: strings.TrimSpace(string(r.Output)),
	})
}
