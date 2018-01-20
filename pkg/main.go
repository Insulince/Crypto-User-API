package main

import (
	"io/ioutil"
	"encoding/json"
	"crypto-users/pkg/models"
	"math/rand"
	"time"
	"crypto-users/pkg/database"
	"log"
	"net/http"
	"strconv"
	"crypto-users/pkg/routes"
)

var config models.Config

func main() () {
	configure()

	router := models.CreateRouter()
	router = routes.CreateRoutes(router)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), router))
}

func configure() () {
	populateConfig()
	rand.Seed(time.Now().UnixNano())
	err := database.InitializeDatabase(config)
	if err != nil {
		panic(err)
	}
	err = database.InsertToken(models.NewToken("is master"))
	if err != nil {
		panic(err)
	}
}

func populateConfig() () {
	jsonFile, err := ioutil.ReadFile("../config/config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		panic(err)
	}
}
