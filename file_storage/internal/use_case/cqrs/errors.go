package cqrs

import (
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"fmt"
)

// var InvalidFileErr = errors.New("Cannot read file from content.")

func SaveFileErr(name string) error {
	return models.NewCommonError(errors.New(fmt.Sprintf("Unable to save file %s", name)))
}

func SavePathErr(path string) error {
	return models.NewCommonError(errors.New(fmt.Sprintf("Unable to save file content to the %s", path)))
}
