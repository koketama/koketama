package util

import (
	"io/ioutil"
	"os"
)

// ReadFile return file content
func ReadFile(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

// WriteFile write some into a file
func WriteFile(name string, raw []byte) error {
	info, err := os.Stat(name)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(name, raw, info.Mode())
}
