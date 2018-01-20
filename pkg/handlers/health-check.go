package handlers

import (
	"net/http"
	"crypto-users/pkg/models/responses"
	"fmt"
	"os"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) () {
	_, _, _, err := CallReceived(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Message{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	Respond(w, responses.Message{Message: "OK"}, http.StatusOK)
}
