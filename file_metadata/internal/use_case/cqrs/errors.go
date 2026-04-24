package cqrs

import (
	"errors"
)

func SaveFileMetadataErr() error {
	return errors.New("Unable to save file metadata")
}
