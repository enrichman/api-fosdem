package store_test

import (
	"fmt"
	"testing"

	"github.com/enrichman/api-fosdem/indexer"
	"gopkg.in/mgo.v2/bson"

	"github.com/enrichman/api-fosdem/store"
)

func TestUnMa(t *testing.T) {
	b, _ := bson.Marshal(indexer.Speaker{
		ID: 123,
	})
	var m map[string]interface{}
	e := bson.Unmarshal(b, m)
	fmt.Println(m, e)
}

func TestGetByID(t *testing.T) {
	s, err := store.NewMongoStore("", "")
	if err != nil {
		panic(err)
	}

	speaker, err := s.FindByID(6)
	fmt.Println(speaker, err)
}
