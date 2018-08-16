package data

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"runtime"
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
	if runtime.GOOS != "windows" {
		clean = strings.Replace(clean, "\\", "", -1)
	}
	if h.hasApostrophes(clean) || h.hasQuotes(clean) {
		clean = clean[1 : len(clean)-1]
	}
	if !h.isGif(clean) {
		err = fmt.Errorf("not a gif [%v]", clean)
	}
	if strings.Count(clean, ".gif") > 1 {
		err = fmt.Errorf("multiple gifs detected in %v", clean)
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

// MD5Checksum returns the checksum for the passed file
func (h Handler) MD5Checksum(filePath string) (string, error) {
	var returnMD5String string

	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}
