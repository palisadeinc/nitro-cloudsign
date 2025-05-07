package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type configHandler struct {
	resp []byte
}

func newConfigHandler(config map[string]string) (*configHandler, error) {
	resp, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal config")
	}

	return &configHandler{resp: resp}, nil
}

func (c *configHandler) Handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.WithError(err).Error("error closing request body")
		}
	}()

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(c.resp)
}
