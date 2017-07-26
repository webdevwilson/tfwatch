package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webdevwilson/tfwatch/test"
)

func TestMain(m *testing.M) {
	test.SuppressLogs()
	os.Exit(m.Run())
}
func Test_systemController(t *testing.T) {
	opts := &Options{
		CheckoutDir: "/opts/repos",
		ClearState:  false,
		LogDir:      "/opts/repos/.tfwatch/logs",
		LogLevel:    "INFO",
		Port:        3000,
		RunPlan:     true,
		SiteDir:     "site",
		StateDir:    "/opts/repos/.tfwatch",
	}

	controller := systemController(opts, nil)
	cfg := controller.GetConfiguration()

	assert.Equal(t, 3, len(cfg))
	return
}
