package models

import (
	"time"
	"math/rand"
	"gopkg.in/mgo.v2/bson"
)

type Token struct {
	Id                bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Value             string        `json:"value" bson:"value"`
	Invalidated       bool          `json:"invalidated" bson:"invalidated"`
	MasterTokenValue  string        `json:"master-token-value" bson:"master-token-value"`
	CreationTimestamp int64         `json:"creation-timestamp" bson:"creation-timestamp"`
}

func NewToken(masterTokenValue string) (token Token) {
	return Token{Id: bson.NewObjectId(), Value: generateRandomTokenValue(), Invalidated: false, MasterTokenValue: masterTokenValue, CreationTimestamp: time.Now().Unix()}
}

func generateRandomTokenValue() (value string) {
	for i := 0; i < 32; i++ {
		value += randomCharacter()
	}

	return value
}

func randomCharacter() (character string) {
	characterPool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`~-_=+[{]}\\|;:.>,</?'\""

	return string(characterPool[rand.Intn(len(characterPool))])
}
