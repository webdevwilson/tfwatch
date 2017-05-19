package routes

import "net/http"

func init() {
	registrationCh <- func(s *server) {
		s.registerAPIEndpoints([]api{
			api{"GET", "/api/configuration", configurationGet},
		}...)
	}
}

func configurationGet(req *http.Request) (data interface{}, err error) {
	return serverSingleton.instance.system.GetConfiguration(), nil
}
