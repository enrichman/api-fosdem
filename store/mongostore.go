package store

import (
	"errors"

	"github.com/enrichman/api-fosdem/indexer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	defaultDB = "api-fosdem"
)

type MongoStore struct {
	db *mgo.Database
}

func NewMongoStore(uri, db string) (*MongoStore, error) {
	if db == "" {
		db = defaultDB
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}
	return &MongoStore{session.DB(db)}, nil
}

func (ms *MongoStore) Save(s indexer.Speaker) error {
	c := ms.db.C("speakers")
	_, err := c.Upsert(bson.M{"_id": s.ID}, s)
	return err
}

func (ms *MongoStore) GetByID(ID int) (*indexer.Speaker, error) {
	c := ms.db.C("speakers")
	iter := c.Find(bson.M{"_id": ID}).Iter()

	var s indexer.Speaker
	if iter.Next(&s) {
		return &s, nil
	}
	return nil, errors.New("not found")
}
