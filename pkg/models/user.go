package models

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id                bson.ObjectId   `json:"_id" bson:"_id,omitempty"`
	Username          string          `json:"username" bson:"username"`
	PasswordHash      []byte          `json:"password-hash" bson:"password-hash"`
	ContactMethods    []ContactMethod `json:"contact-methods" bson:"contact-methods"`
	WatchList         []string        `json:"watch-list" bson:"watch-list"`
	CreationTimestamp int64           `json:"creation-timestamp" bson:"creation-timestamp"`
}

type ContactMethod struct {
	Type  string `json:"type" bson:"type"`
	Value string `json:"value" bson:"value"`
}
