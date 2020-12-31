package utility

import (
	"errors"
	"io/ioutil"
	"os"
)

// ReadJSON Reads data from JSON onto a struct
func ReadJSON(filepath string) ([]byte, error) {
	if jsonFile, err := os.Open(filepath); err == nil {
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		return byteValue, nil
	} else {
		return []byte{}, errors.New("could not read config file")
	}
}