package handlers

import (
	"net/http"
	"fmt"
	"os"
	"crypto-users/pkg/models/responses"
)

func NotFound(w http.ResponseWriter, r *http.Request) () {
	_, _, _, err := CallReceived(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Message{Message: err.Error()})
		return
	}

	Respond(w, responses.Message{Message: "Unrecognized call."})
}
