package models

type Config struct {
	Port        int    `json:"port"`
	MongoDBURL  string `json:"mongo-db-url"`
}
