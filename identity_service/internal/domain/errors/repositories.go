package errors

import (
	"elex_storage/pkg/shared_kernel/models"
	"errors"
)

var ErrInsertUser = models.NewCommonError(errors.New("Can't create new user"))
var ErrUserExists = models.NewCommonError(errors.New("Same user is already exists"))
