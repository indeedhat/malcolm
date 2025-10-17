package malcolm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var prefixTestCases = []struct {
	base     string
	url      string
	expected string
}{
	{"", "", ""},
	{"", "/", "/"},
	{"", "/path", "/path"},
	{"", "/path/multi/part", "/path/multi/part"},
	{"", "/path/:arg", "/path/:arg"},
	{"", "GET /path/:arg", "GET /path/:arg"},

	{"/", "", "/"},
	{"/", "/", "/"},
	{"/", "/path", "/path"},
	{"/", "/path/multi/part", "/path/multi/part"},
	{"/", "/path/:arg", "/path/:arg"},
	{"/", "GET /path/:arg", "GET /path/:arg"},

	{"/base", "", "/base"},
	{"/base", "/", "/base"},
	{"/base", "/path", "/base/path"},
	{"/base", "/path/multi/part", "/base/path/multi/part"},
	{"/base", "/path/:arg", "/base/path/:arg"},
	{"/base", "GET /path/:arg", "GET /base/path/:arg"},

	{"/multi/base", "", "/multi/base"},
	{"/multi/base", "/", "/multi/base"},
	{"/multi/base", "/path", "/multi/base/path"},
	{"/multi/base", "/path/multi/part", "/multi/base/path/multi/part"},
	{"/multi/base", "/path/:arg", "/multi/base/path/:arg"},
	{"/multi/base", "GET /path/:arg", "GET /multi/base/path/:arg"},
}

func TestPrefix(t *testing.T) {
	r := NewDefaultRouter()

	for _, tcase := range prefixTestCases {
		t.Run(fmt.Sprint(tcase.base, " -> ", tcase.url), func(t *testing.T) {
			g := r.Group(tcase.base)
			final := g.prefix(tcase.url)

			assert.Equal(t, tcase.expected, final)
		})
	}
}

func TestGroup(t *testing.T) {
	r := NewDefaultRouter()
	g := r.Group("/base/path", func(next http.HandlerFunc) http.HandlerFunc {
		return next
	})

	assert.Equal(t, r.mux, g.mux)
	assert.Len(t, g.middleware, 1)
	assert.Equal(t, "/base/path", g.basePath)

	g2 := g.Group("/extended", func(next http.HandlerFunc) http.HandlerFunc {
		return next
	})

	assert.Equal(t, r.mux, g2.mux)
	assert.Len(t, g2.middleware, 2)
	assert.Equal(t, "/base/path/extended", g2.basePath)
}
