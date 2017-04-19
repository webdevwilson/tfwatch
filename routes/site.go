package routes

import "github.com/webdevwilson/terraform-ci/config"

func init() {
	//r := Router()
	_ = config.Get()

	// Serve /site directory
	//fs := http.FileServer(http.Dir(cfg.SiteRoot))
	//r.PathPrefix("/static").Handler(http.StripPrefix("/", fs))
}
