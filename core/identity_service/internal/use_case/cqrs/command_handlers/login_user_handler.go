package command_handlers

import (
	"elex_storage/identity_service/internal/core_utils"
	"elex_storage/identity_service/internal/domain/entities"
	"elex_storage/identity_service/internal/domain/messaging/publishers"
	"elex_storage/identity_service/internal/domain/repositories"
	"elex_storage/identity_service/internal/use_case/cqrs"
	"elex_storage/identity_service/internal/use_case/cqrs/commands"
	"elex_storage/identity_service/internal/use_case/cqrs/responses"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/input_validation"
	shared_models "elex_storage/pkg/shared_kernel/models"
	"elex_storage/pkg/shared_kernel/token_handlers"
	"net/mail"
)

type LoginUserHandler struct {
	logger           logger.Logger
	userRepository   repositories.UserRepository
	storagePublisher publishers.StoragePublisher
	identityUtil     *core_utils.IdentityUtil
	jwtToken         *token_handlers.JwtToken
}

func NewLoginUserHandler(deps cqrs.UserHandlerDeps) *LoginUserHandler {
	return &LoginUserHandler{
		logger:           deps.Logger,
		userRepository:   deps.UserRepository,
		storagePublisher: deps.StoragePublisher,
		identityUtil:     deps.IdentityUtil,
		jwtToken:         deps.JwtToken,
	}
}

func (handler *LoginUserHandler) Handle(cmd commands.LoginUserCommand) (*responses.TokenResponse, error) {
	// First step: check the UserName is phone or email
	var user *entities.UserEntity
	if _, err := mail.ParseAddress(cmd.UserName); err == nil {
		user, err = handler.userRepository.GetUserByEmail(cmd.UserName)
		if err != nil {
			return nil, err
		}
	} else if err = input_validation.IsPhoneNumber(cmd.UserName); err == nil {
		phoneInfo, err := input_validation.CheckPhone(cmd.UserName)
		if err != nil {
			cErr := shared_models.NewCommonError(err)
			return nil, cErr
		}
		user, err = handler.userRepository.GetUserByPhone(phoneInfo.FixPhoneNumber)
	}
	if user == nil {
		return nil, cqrs.UserNotExistErr
	}
	// Second step: Compare Hash the password
	if !handler.identityUtil.CheckPasswordHash(cmd.Password, user.PasswordHash) {
		return nil, cqrs.InvalidPasswordErr
	}
	tokenPair, err := handler.jwtToken.GenerateTokenPair(user.Id.String(), "")
	if err != nil {
		return nil, generateTokenErr
	}
	response := responses.TokenResponse{
		UserId:       user.Id,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
	return &response, nil
}
