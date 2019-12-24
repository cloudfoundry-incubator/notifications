package info

import (
	"encoding/json"
	"net/http"

	"github.com/ryanmoran/stack"
)

type GetHandler struct {
}

func NewGetHandler() GetHandler {
	return GetHandler{}
}

func (handler GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	output, err := json.Marshal(map[string]interface{}{
		"version": 1,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
