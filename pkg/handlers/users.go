package handlers

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
	"crypto-users/pkg/models"
	"crypto-users/pkg/database"
	"time"
	"fmt"
	"os"
	"crypto-users/pkg/models/responses"
	"crypto-users/pkg/services"
)

func Register(w http.ResponseWriter, r *http.Request) () {
	_, _, rawPostBody, err := CallReceived(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not process request."}, http.StatusInternalServerError)
		return
	}

	type PostBody struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm-password"`
	}
	var postBody PostBody
	err = json.Unmarshal(rawPostBody, &postBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not read request body."}, http.StatusBadRequest)
		return
	}

	if postBody.Password != postBody.ConfirmPassword {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Passwords did not match."}, http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(postBody.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not generate password hash."}, http.StatusInternalServerError)
		return
	}

	err = database.InsertUser(models.User{Username: postBody.Username, PasswordHash: passwordHash, ContactMethods: []models.ContactMethod{}, WatchList: []string{}, CreationTimestamp: time.Now().Unix()})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not insert user."}, http.StatusInternalServerError)
		return
	}

	Respond(w, responses.Message{Message: "User registered successfully."}, http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) () {
	_, _, rawPostBody, err := CallReceived(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not process request."}, http.StatusInternalServerError)
		return
	}

	type PostBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var postBody PostBody
	err = json.Unmarshal(rawPostBody, &postBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not read request body."}, http.StatusBadRequest)
		return
	}

	user, err := database.FindUserByUsername(postBody.Username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Invalid username or password."}, http.StatusBadRequest)
		return
	}
	usersPasswordHash := user.PasswordHash
	err = bcrypt.CompareHashAndPassword(usersPasswordHash, []byte(postBody.Password))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Invalid username or password."}, http.StatusBadRequest)
		return
	}

	masterToken, err := database.GetMasterToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not locate current master token."}, http.StatusInternalServerError)
		return
	}
	token := models.NewToken(masterToken.Value)
	err = database.InsertToken(token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not insert token."}, http.StatusInternalServerError)
		return
	}

	Respond(w, responses.Token{Message: "Success.", Token: token}, http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) () {
	_, queryParameters, _, err := CallReceived(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not process request."}, http.StatusInternalServerError)
		return
	}

	if len(queryParameters["token-id"]) != 1 {
		fmt.Fprintln(os.Stderr, "No token-id")
		Respond(w, responses.Error{Message: "No \"token-id\" query parameter provided for authentication, access denied."}, http.StatusInternalServerError)
		return
	}

	if len(queryParameters["token-value"]) != 1 {
		fmt.Fprintln(os.Stderr, "No token-value provided.")
		Respond(w, responses.Error{Message: "No \"token-value\" query parameter was provided, access denied."}, http.StatusBadRequest)
		return
	}

	valid, message := services.VerifyToken(queryParameters["token-id"][0], queryParameters["token-value"][0])
	if valid != true {
		fmt.Fprintln(os.Stderr, message)
		Respond(w, responses.Error{Message: "Could not invalidate token: " + message}, http.StatusInternalServerError)
		return
	}

	err = database.InvalidateTokenWithId(queryParameters["token-id"][0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not invalidate token."}, http.StatusInternalServerError)
		return
	}
	Respond(w, responses.Message{Message: "Success."}, http.StatusOK)
}

func escape(qs string) (escapedQS string) {
	escapedQS = ""

	for i := 0; i < len(qs); i++ {
		escapedQS += string(qs[i])

		if qs[i] == '%' {
			escapedQS += "%"
		}
	}

	return escapedQS
}
