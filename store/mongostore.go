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

// MongoStore can save and retrieve Speakers from MongoDB
type MongoStore struct {
	db *mgo.Database
}

// NewMongoStore creates a new MongoStore
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

// Save a speaker of the passed year
func (ms *MongoStore) Save(s indexer.Speaker, year int) error {
	c := ms.db.C(speakerCollection)
	_, err := c.Upsert(
		bson.M{"_id": s.ID},
		bson.M{
			"$set": bson.M{
				"slug":         s.Slug,
				"name":         s.Name,
				"profileimage": s.ProfileImage,
				"profilepage":  s.ProfilePage,
				"links":        s.Links,
				"bio":          s.Bio,
			},
			"$addToSet": bson.M{"years": year},
		},
	)
	return err
}

// FindByID find a Speaker from its ID
func (ms *MongoStore) FindByID(ID int) (*indexer.Speaker, error) {
	c := ms.db.C(speakerCollection)
	iter := c.Find(bson.M{"_id": ID}).Iter()

	var s indexer.Speaker
	if iter.Next(&s) {
		return &s, nil
	}
	return nil, errors.New("not found")
}

// Find find a list of Speakers based on the passed params
func (ms *MongoStore) Find(limit, offset int, slug string, years []int) ([]indexer.Speaker, int, error) {
	c := ms.db.C(speakerCollection)

	ors := make([]bson.M, 0)
	for _, n := range strings.Split(slug, " ") {
		ors = append(ors, bson.M{"slug": bson.RegEx{Pattern: n, Options: "i"}})
	}
	if len(years) > 0 {
		ors = append(ors, bson.M{"years": bson.M{"$in": years}})
	}

	query := c.Find(bson.M{"$and": ors})

	count, err := query.Count()
	if err != nil {
		return nil, 0, err
	}

	iter := query.Skip(offset).Limit(limit).Sort("_id").Iter()
	speakersFound := make([]indexer.Speaker, 0)
	var s indexer.Speaker
	for iter.Next(&s) {
		speakersFound = append(speakersFound, s)
	}

	return speakersFound, count, nil
}
