package errors

import (
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"fmt"
)

var ErrInsertFileMetadata = models.NewCommonError(errors.New("Can't insert file metadata into db"))

func ErrDirectoryExists(dirName string) error {
	return models.NewCommonError(errors.New(
		fmt.Sprintf("Directory with the name %s already exists", dirName)))
}

var ErrFileExists = models.NewCommonError(errors.New("File with the same name already exists"))
