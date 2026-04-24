package cqrs

import (
	"elex_storage/identity_service/internal/core_utils"
	"elex_storage/identity_service/internal/domain/messaging/publishers"
	"elex_storage/identity_service/internal/domain/repositories"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/models"
	"elex_storage/pkg/shared_kernel/token_handlers"
	"errors"

	"go.uber.org/fx"
)

type UserHandlerDeps struct {
	fx.In

	Logger           logger.Logger
	UserRepository   repositories.UserRepository
	StoragePublisher publishers.StoragePublisher
	IdentityUtil     *core_utils.IdentityUtil
	JwtToken         *token_handlers.JwtToken
}

var GenerateTokenErr = models.NewCommonError(errors.New("Can't create new token"))
var UserNotExistErr = models.NewCommonError(errors.New("User is not available"))
var InvalidPasswordErr = models.NewCommonError(errors.New("Password is not correct"))
