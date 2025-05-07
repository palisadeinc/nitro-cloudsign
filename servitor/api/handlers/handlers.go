package handlers

import "net/http"

func Handler(config map[string]string) (http.Handler, error) {
	cHandler, err := newConfigHandler(config)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/config", cHandler.Handler)
	return mux, nil
}
