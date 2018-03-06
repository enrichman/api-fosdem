package speakers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeSpeakersHandler setup the handlers on the /api/v1/speakers route
func MakeSpeakersHandler(s speakerService) http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	speakerGetterHandler := kithttp.NewServer(
		makeSpeakerGetterEndpoint(s),
		decodeSpeakerGetter,
		encodeSpeakerGetter,
	)

	speakerFinderHandler := kithttp.NewServer(
		makeSpeakerFinderEndpoint(s),
		decodeSpeakerFinder,
		encodeSpeakerFinder,
	)

	r.Handle("/api/v1/speakers", speakerFinderHandler).Methods(http.MethodGet)
	r.Handle("/api/v1/speakers/{id}", speakerGetterHandler).Methods(http.MethodGet)

	return r
}

type getSpeakerByIDRequest struct {
	id   int
	year int
}

func decodeSpeakerGetter(_ context.Context, r *http.Request) (interface{}, error) {
	var err error
	req := getSpeakerByIDRequest{year: 2018}

	vars := mux.Vars(r)
	req.id, err = strconv.Atoi(vars["id"])
	if err != nil {
		return nil, errors.New("wrong ID")
	}

	if yearStr := r.FormValue("year"); yearStr != "" {
		req.year, err = strconv.Atoi(yearStr)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

func encodeSpeakerGetter(_ context.Context, w http.ResponseWriter, res interface{}) error {
	return json.NewEncoder(w).Encode(res)
}

func decodeSpeakerFinder(_ context.Context, r *http.Request) (interface{}, error) {
	var err error
	var req findRequest

	if limit := r.FormValue("limit"); limit != "" {
		req.limit, err = strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
	}

	if offset := r.FormValue("offset"); offset != "" {
		req.offset, err = strconv.Atoi(offset)
		if err != nil {
			return nil, err
		}
	}

	req.years = make([]int, 0)
	if year := r.FormValue("year"); year != "" {
		for _, y := range strings.Split(year, ",") {
			yearInt, err := strconv.Atoi(y)
			if err != nil {
				return nil, err
			}
			req.years = append(req.years, yearInt)
		}
	}
	if len(req.years) == 0 {
		req.years = append(req.years, 2018)
	}

	req.slug = r.FormValue("slug")

	if req.limit == 0 || req.limit > 100 {
		req.limit = 100
	}
	return req, nil
}

func encodeSpeakerFinder(_ context.Context, w http.ResponseWriter, res interface{}) error {
	return json.NewEncoder(w).Encode(res)
}
