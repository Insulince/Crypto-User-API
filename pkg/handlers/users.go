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
)

func Register(w http.ResponseWriter, r *http.Request) () {
	_, _, rawPostBody, err := CallReceived(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not process request."})
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
		Respond(w, responses.Error{Message: "Could not read request body."})
		return
	}

	if postBody.Password != postBody.ConfirmPassword {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Passwords did not match."})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(postBody.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not generate password hash."})
		return
	}

	err = database.InsertUser(models.User{Username: postBody.Username, PasswordHash: passwordHash, ContactMethods: []models.ContactMethod{}, WatchList: []string{}, CreationTimestamp: time.Now().Unix()})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not insert user."})
		return
	}

	Respond(w, responses.Message{Message: "User registered successfully."})
}

func Login(w http.ResponseWriter, r *http.Request) () {
	_, _, rawPostBody, err := CallReceived(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not process request."})
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
		Respond(w, responses.Error{Message: "Could not read request body."})
		return
	}

	user, err := database.FindUserByUsername(postBody.Username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Invalid username or password."})
		return
	}
	usersPasswordHash := user.PasswordHash
	err = bcrypt.CompareHashAndPassword(usersPasswordHash, []byte(postBody.Password))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Invalid username or password."})
		return
	}

	masterToken, err := database.GetMasterToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not locate current master token."})
		return
	}
	token := models.NewToken(masterToken.Value)
	err = database.InsertToken(token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not insert token."})
		return
	}
	type Response struct {
		Message string       `json:"message"`
		Token   models.Token `json:"token"`
	}
	Respond(w, Response{Message: "Success.", Token: token})
}

func Logout(w http.ResponseWriter, r *http.Request) () {
	_, _, rawPostBody, err := CallReceived(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not process request."})
		return
	}

	type PostBody struct {
		TokenValue string `json:"token-value"`
	}
	var postBody PostBody
	err = json.Unmarshal(rawPostBody, &postBody)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not read request body."})
		return
	}

	err = database.InvalidateTokenWithValue(postBody.TokenValue)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		Respond(w, responses.Error{Message: "Could not invalidate token."})
		return
	}
	Respond(w, responses.Message{Message: "Success."})
}
