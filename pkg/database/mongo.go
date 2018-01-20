package database

import (
	"gopkg.in/mgo.v2"
	"crypto-users/pkg/models"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

var db *mgo.Database

func InitializeDatabase(config models.Config) (err error) {
	session, err := mgo.Dial(config.MongoDBURL)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Strong, true)
	db = session.DB("crypto")
	return nil
}

func Users() (users *mgo.Collection) {
	return db.C("users")
}

func InsertUser(user models.User) (err error) {
	return Users().Insert(user)
}

func FindUsers() (users []models.User, err error) {
	err = Users().Find(nil).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func FindUserById(id string) (user *models.User, err error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Provided ID \"" + id + "\" is not a valid MongoDB ID.")
	}

	err = Users().FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindUserByUsername(username string) (user *models.User, err error) {
	err = Users().Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(id string, updates bson.M) (err error) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Provided ID \"" + id + "\" is not a valid MongoDB ID.")
	}
	return Users().UpdateId(bson.ObjectIdHex(id), bson.M{"$set": updates})
}

func DeleteUser(id string) (err error) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Provided ID \"" + id + "\" is not a valid MongoDB ID.")
	}
	return Users().RemoveId(bson.ObjectIdHex(id))
}

func Tokens() (tokens *mgo.Collection) {
	return db.C("tokens")
}

func InsertToken(token models.Token) (err error) {
	return Tokens().Insert(token)
}

func FindTokens() (tokens []models.Token, err error) {
	err = Tokens().Find(nil).All(&tokens)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func FindTokenById(id string) (token *models.Token, err error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Provided ID \"" + id + "\" is not a valid MongoDB ID.")
	}

	err = Tokens().FindId(bson.ObjectIdHex(id)).One(&token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func FindTokenByValue(value string) (token *models.Token, err error) {
	err = Tokens().Find(bson.M{"value": value}).One(&token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func UpdateToken(id string, updates bson.M) (err error) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Provided ID \"" + id + "\" is not a valid MongoDB ID.")
	}
	return Tokens().UpdateId(bson.ObjectIdHex(id), bson.M{"$set": updates})
}

func DeleteToken(id string) (err error) {
	if !bson.IsObjectIdHex(id) {
		return errors.New("Provided ID \"" + id + "\" is not a valid MongoDB ID.")
	}
	return Tokens().RemoveId(bson.ObjectIdHex(id))
}

func InvalidateTokenWithId(id string) (err error) {
	token, err := FindTokenById(id)
	if err != nil {
		return err
	}
	return UpdateToken(token.Id.Hex(), bson.M{"invalidated": true})
}

func InvalidateTokenWithValue(value string) (err error) {
	token, err := FindTokenByValue(value)
	if err != nil {
		return err
	}
	return UpdateToken(token.Id.Hex(), bson.M{"invalidated": true})
}

func GetMasterToken() (masterToken *models.Token, err error) {
	err = Tokens().Find(bson.M{"master-token-value": "is master"}).Sort("-creation-timestamp").One(&masterToken)
	if err != nil {
		return nil, err
	}
	return masterToken, nil
}
