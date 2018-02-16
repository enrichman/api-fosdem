package store

import (
	"errors"
	"strings"

	"github.com/enrichman/api-fosdem/indexer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	defaultDB         = "api-fosdem"
	speakerCollection = "speakers"
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
	c := ms.db.C(speakerCollection)
	_, err := c.Upsert(bson.M{"_id": s.ID}, s)
	return err
}

func (ms *MongoStore) FindByID(ID int) (*indexer.Speaker, error) {
	c := ms.db.C(speakerCollection)
	iter := c.Find(bson.M{"_id": ID}).Iter()

	var s indexer.Speaker
	if iter.Next(&s) {
		return &s, nil
	}
	return nil, errors.New("not found")
}

func (ms *MongoStore) Find(name string) []indexer.Speaker {
	c := ms.db.C(speakerCollection)

	ors := make([]bson.M, 0)
	for _, n := range strings.Split(name, " ") {
		ors = append(ors, bson.M{"slug": bson.RegEx{Pattern: n, Options: "i"}})
	}
	query := bson.M{"$and": ors}
	iter := c.Find(query).Iter()

	speakersFound := make([]indexer.Speaker, 0)

	var s indexer.Speaker
	for iter.Next(&s) {
		speakersFound = append(speakersFound, s)
	}
	return speakersFound
}
