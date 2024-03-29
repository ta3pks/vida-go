package models

import (
	"fmt"
	"time"

	"gitlab.mugsoft.io/vida/go-api/config"
	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Database

const DATA_PER_PAGE = 10

func db_get() *mgo.Database {
	if nil != db {
		return db
	}
	fmt.Println("trying to connect mongo")
	_ses, err := mgo.DialWithTimeout(config.Get("DB_ADDR"), time.Second*2)
	if nil != err {
		panic(err)
	}
	db = _ses.DB(config.Get("DB"))
	return db
}

type Defaults struct {
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func _update_fields(q map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{"$set": q}
}
