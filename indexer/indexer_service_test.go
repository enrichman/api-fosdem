package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getSlugByLink(t *testing.T) {
	tt := []struct {
		name string
		link string
		slug string
	}{
		{
			name: "normal link",
			link: "/2018/schedule/speaker/a_simple_slug/",
			slug: "a_simple_slug",
		},
		{
			name: "normal link with no ending slash",
			link: "/2018/schedule/speaker/a_simple_slug",
			slug: "a_simple_slug",
		},
		{
			name: "simple string",
			link: "a_simple_slug",
			slug: "a_simple_slug",
		},
		{
			name: "empty string",
			link: "",
			slug: "",
		},
		{
			name: "only slash",
			link: "/",
			slug: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			slug := getSlugByLink(tc.link)
			assert.Equal(t, tc.slug, slug)
		})
	}
}
