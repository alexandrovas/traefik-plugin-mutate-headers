// nolint
package traefik_plugin_mutate_headers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

// Mutation holds one mutation body configuration.
type Mutation struct {
	Header       string `json:"header,omitempty"`
	NewName      string `json:"newName,omitempty"`
	DeleteSource bool   `json:"deleteSource,omitempty"`
	Regex        string `json:"regex,omitempty"`
	Replacement  string `json:"replacement,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	Mutations []Mutation `json:"mutations,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type mutation struct {
	oldName      string
	newName      string
	deleteSource bool
	mutate       bool
	regex        *regexp.Regexp
	replacement  string
}

type HeaderMutator struct {
	name      string
	next      http.Handler
	mutations []mutation
}

// New creates and returns a new HeaderMutator plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	mutations := make([]mutation, len(config.Mutations))

	for i, m := range config.Mutations {
		mt := mutation{oldName: m.Header, newName: m.NewName, deleteSource: m.DeleteSource}
		if m.Regex != "" {
			regex, err := regexp.Compile(m.Regex)
			if err != nil {
				return nil, fmt.Errorf("error compiling regex %q: %w", m.Regex, err)
			}
			if m.Replacement == "" {
				return nil, fmt.Errorf("replacement is required when regex is set")
			}
			mt.mutate = true
			mt.regex = regex
			mt.replacement = m.Replacement
		} else {
			mt.mutate = false
		}

		if m.NewName == "" {
			mt.newName = m.Header
		}

		mutations[i] = mt
	}

	return &HeaderMutator{
		name:      name,
		next:      next,
		mutations: mutations,
	}, nil
}

func (h *HeaderMutator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, m := range h.mutations {
		headerValues := req.Header.Values(m.oldName)

		if m.deleteSource {
			req.Header.Del(m.oldName)
		}

		if len(headerValues) == 0 {
			continue
		}

		newHeader := req.Header.Get(m.newName)
		if newHeader != "" {
			req.Header.Del(m.newName)
		}

		for _, v := range headerValues {
			if m.mutate {
				mv := m.regex.ReplaceAllString(v, m.replacement)
				if mv != "" {
					req.Header.Add(m.newName, mv)
				} else {
					req.Header.Add(m.newName, v)
				}
			} else {
				req.Header.Add(m.newName, v)
			}
		}
	}

	h.next.ServeHTTP(rw, req)
}
