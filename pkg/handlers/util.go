package handlers

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
)

func CallReceived(r *http.Request) (routeVariables map[string]string, queryParameters map[string][]string, requestBody []byte, err error) {
	fmt.Printf("Call Received: \"" + r.Method + " " + r.URL.Path + "\"\n")
	return getRequestInformation(r)
}

func getRequestInformation(r *http.Request) (routeVariables map[string]string, queryParameters map[string][]string, requestBody []byte, err error) {
	requestBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, nil, nil, err
	}
	return mux.Vars(r), r.URL.Query(), requestBody, nil
}

func Respond(w http.ResponseWriter, response interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json")
	responseBody, err := json.Marshal(response)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, string(responseBody))
	return nil
}