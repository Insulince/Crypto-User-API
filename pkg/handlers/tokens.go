package handlers

import (
	"net/http"
	"fmt"
	"os"
	"crypto-users/pkg/models/responses"
	"crypto-users/pkg/services"
)

func Verify(w http.ResponseWriter, r *http.Request) () {
	_, queryParameters, _, err := CallReceived(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not process request."}, http.StatusInternalServerError)
		return
	}

	if len(queryParameters["token-id"]) != 1 {
		fmt.Fprintln(os.Stderr, "No token-id provided.")
		Respond(w, responses.Error{Message: "No \"token-id\" query parameter was provided, cannot verify."}, http.StatusBadRequest)
		return
	}

	if len(queryParameters["token-value"]) != 1 {
		fmt.Fprintln(os.Stderr, "No token-value provided.")
		Respond(w, responses.Error{Message: "No \"token-value\" query parameter was provided, cannot verify."}, http.StatusBadRequest)
		return
	}

	valid, message := services.VerifyToken(queryParameters["token-id"][0], queryParameters["token-value"][0])

	Respond(w, responses.StatusMessage{Valid: valid, Message: message}, http.StatusOK)
}
