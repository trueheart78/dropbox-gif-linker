package data

import (
	"fmt"
	"strings"
)

// Handler creates a new input handler
type Handler struct {
}

// NewHandler returns a new Handler
func NewHandler() Handler {
	return Handler{}
}

// Clean cleans the input data
func (h Handler) Clean(data string) (clean string, err error) {
	clean = strings.TrimSpace(data)
	clean = strings.Replace(clean, "\\", "", -1)
	if h.hasApostrophes(clean) || h.hasQuotes(clean) {
		clean = clean[1 : len(clean)-1]
	}
	if !h.isGif(clean) {
		err = fmt.Errorf("not a gif [%v]", clean)
	}
	return
}

func (h Handler) isGif(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".gif")
}

func (h Handler) hasApostrophes(data string) bool {
	return (strings.HasPrefix(data, "'") && strings.HasSuffix(data, "'"))
}

func (h Handler) hasQuotes(data string) bool {
	return (strings.HasPrefix(data, "\"") && strings.HasSuffix(data, "\""))

}
