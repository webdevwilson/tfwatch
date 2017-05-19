package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_systemController(t *testing.T) {
	opts := &Options{
		CheckoutDir: "/opts/repos",
		ClearState:  false,
		LogDir:      "/opts/repos/.terraform-ci/logs",
		LogLevel:    "INFO",
		Port:        3000,
		RunPlan:     true,
		SiteDir:     "site",
		StateDir:    "/opts/repos/.terraform-ci",
	}

	controller := systemController(opts, nil)
	cfg := controller.GetConfiguration()

	assert.Equal(t, 3, len(cfg))
	return
}
