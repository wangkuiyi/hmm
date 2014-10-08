package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

// LoadModel decodes JSON-encoded model file filename.
func LoadModel(filename string) (*Model, error) {
	m := new(Model)

	if len(filename) <= 0 {
		return nil, errors.New("Cannot load model")
	}

	if f, e := os.Open(filename); e != nil {
		return nil, e
	} else {
		defer f.Close()
		if e := json.NewDecoder(f).Decode(m); e != nil {
			return nil, e
		}
	}

	return m, nil
}

// SaveModel writes JSON-encoded model to the file named by filename.
// If filename is an empty string, it writes to standard output.
func SaveModel(m *Model, filename string) {
	f := CreateFileOrStdout(filename)
	if f != os.Stdout {
		defer f.Close()
	}

	if b, e := json.MarshalIndent(m, "", "  "); e != nil {
		log.Fatalf("Failed encoding model: %v", e)
	} else {
		fmt.Fprintf(f, "%s", b)
	}
}

// CreateFileOrStdout creates a new file named by filename, or returns
// simply standard output if filename is an empty string.  Be aware to
// close the created file but not to close standard output.
func CreateFileOrStdout(filename string) io.WriteCloser {
	if f, e := os.Create(filename); e == nil {
		return f
	}
	return os.Stdout
}
