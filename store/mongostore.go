package store

import (
	"errors"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	defaultDB         = "api-fosdem"
	speakerCollection = "speakers"
)

// Speaker maps the speaker
type Speaker struct {
	ID           int
	Slug         string
	Name         string
	ProfileImage string
	ProfilePage  string
	Bio          string
	Year         int
	Links        []Link
}

// Link is a detail link owned by a Speaker
type Link struct {
	URL   string
	Title string
}

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
func (ms *MongoStore) Save(s Speaker) error {
	c := ms.db.C(speakerCollection)
	_, err := c.Upsert(
		bson.M{"id": s.ID, "year": s.Year},
		bson.M{
			"$set": bson.M{
				"id":           s.ID,
				"slug":         s.Slug,
				"name":         s.Name,
				"profileimage": s.ProfileImage,
				"profilepage":  s.ProfilePage,
				"bio":          s.Bio,
				"year":         s.Year,
				"links":        s.Links,
			},
		},
	)
	return err
}

// FindByID find a Speaker from its ID
func (ms *MongoStore) FindByID(ID, year int) (*Speaker, error) {
	c := ms.db.C(speakerCollection)
	iter := c.Find(bson.M{"id": ID, "year": year}).Iter()

	var s Speaker
	if iter.Next(&s) {
		return &s, nil
	}
	return nil, errors.New("not found")
}

// Find find a list of Speakers based on the passed params
func (ms *MongoStore) Find(limit, offset int, slug string, years []int) ([]Speaker, int, error) {
	c := ms.db.C(speakerCollection)

	ors := make([]bson.M, 0)
	for _, n := range strings.Split(slug, " ") {
		ors = append(ors, bson.M{"slug": bson.RegEx{Pattern: n, Options: "i"}})
	}
	if len(years) > 0 {
		ors = append(ors, bson.M{"year": bson.M{"$in": years}})
	}

	query := c.Find(bson.M{"$and": ors})

	count, err := query.Count()
	if err != nil {
		return nil, 0, err
	}

	iter := query.Skip(offset).Limit(limit).Sort("id", "-year").Iter()
	speakersFound := make([]Speaker, 0)
	var s Speaker
	for iter.Next(&s) {
		speakersFound = append(speakersFound, s)
	}

	return speakersFound, count, nil
}
