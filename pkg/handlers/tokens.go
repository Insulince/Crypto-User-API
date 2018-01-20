package handlers

import (
	"net/http"
	"crypto-users/pkg/database"
	"fmt"
	"os"
	"crypto-users/pkg/models/responses"
)

func Verify(w http.ResponseWriter, r *http.Request) () {
	_, queryParameters, _, err := CallReceived(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not process request."}, http.StatusInternalServerError)
		return
	}

	if len(queryParameters["token-id"]) != 1 {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "No \"token-id\" query parameter was provided, cannot verify."}, http.StatusBadRequest)
		return
	}

	token, err := database.FindTokenById(queryParameters["token-id"][0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.StatusError{Status: "Invalid", Message: "A token with the provided token id could not be found."}, http.StatusOK)
		return
	}

	if token.Invalidated != false {
		fmt.Fprintln(os.Stderr, "Token is already invalidated.")
		Respond(w, responses.StatusError{Status: "Invalid", Message: "This token has been invalidated, please fetch a new token (login again)."}, http.StatusOK)
		return
	}

	masterToken, err := database.GetMasterToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not locate current master token."}, http.StatusInternalServerError)
		return
	}

	if token.MasterTokenValue != masterToken.Value {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.StatusError{Status: "Invalid", Message: "Master token associated with provided token is invalid, please fetch a new token (login again)."}, http.StatusOK)
		return
	}

	Respond(w, responses.StatusMessage{Status: "Valid", Message: "Provided token is valid."}, http.StatusOK)
}
