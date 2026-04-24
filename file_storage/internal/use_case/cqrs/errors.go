package cqrs

import (
	"errors"
	"fmt"
)

var InvalidFileErr = errors.New("Cannot read file from content.")

func SaveFileErr(path string) error {
	return errors.New(fmt.Sprintf("Unable to save file content to the %s", path))
}
