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
func (h Handler) Clean(data string) string {
	return fmt.Sprintf("Cleaning [%v]", data)
}

func (h Handler) isGif(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".gif")
}
