package store_test

import (
	"fmt"
	"testing"

	"github.com/enrichman/api-fosdem/store"
)

func TestGetByID(t *testing.T) {
	s, err := store.NewMongoStore("", "")
	if err != nil {
		panic(err)
	}

	speaker, err := s.GetByID(6)
	fmt.Println(speaker, err)
}
