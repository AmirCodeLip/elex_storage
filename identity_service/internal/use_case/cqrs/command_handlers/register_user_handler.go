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
	"elex_storage/pkg/shared_kernel/token_handlers"
	"errors"
	"net/mail"

	"elex_storage/pkg/shared_kernel/event_models"
	"elex_storage/pkg/shared_kernel/guard"
	"elex_storage/pkg/shared_kernel/models"
	"time"

	"github.com/google/uuid"
)

type RegisterUserHandler struct {
	logger           logger.Logger
	userRepository   repositories.UserRepository
	storagePublisher publishers.StoragePublisher
	identityUtil     *core_utils.IdentityUtil
	jwtToken         *token_handlers.JwtToken
}

func NewRegisterUserHandler(deps cqrs.UserHandlerDeps) *RegisterUserHandler {
	return &RegisterUserHandler{
		logger:           deps.Logger,
		userRepository:   deps.UserRepository,
		storagePublisher: deps.StoragePublisher,
		identityUtil:     deps.IdentityUtil,
		jwtToken:         deps.JwtToken,
	}
}

var generateTokenErr = models.NewCommonError(errors.New("Can't make token"))
var emailNotValid = models.NewCommonError(errors.New("Provide valid email address"))
var phoneNotValid = models.NewCommonError(errors.New("Provide valid phone number"))
var requiredInput = models.NewCommonError(errors.New("At least one phone or email is required"))

func (handler *RegisterUserHandler) Handle(cmd commands.RegisterUserCommand) (*responses.TokenResponse, error) {
	// Make sure at least one of email or phone number has a value
	if !guard.AgainstPNullStrs(cmd.Email, cmd.Phone) {
		return nil, requiredInput
	}
	if guard.AgainstPNullStr(cmd.Email) {
		// If email has value then validate the value
		_, err := mail.ParseAddress(*cmd.Email)
		if err != nil {
			return nil, emailNotValid
		}
	}
	if guard.AgainstPNullStr(cmd.Phone) {
		err := input_validation.IsPhoneNumber(*cmd.Phone)
		if err != nil {
			return nil, phoneNotValid
		}
	}
	newUserId := uuid.New()
	// First step: Hash the password
	hashedPassword, _ := handler.identityUtil.HashPassword(cmd.Password)
	// Second step: Insert the user into the database.
	userEntity := entities.UserEntity{
		Id:           newUserId,
		Email:        cmd.Email,
		Phone:        cmd.Phone,
		CreatedAt:    time.Now(),
		PasswordHash: hashedPassword,
	}
	insertErr := handler.userRepository.Insert(userEntity)
	if insertErr != nil {
		return nil, insertErr
	}
	handler.storagePublisher.PublishUserRegisterd(event_models.UserRegisterd{
		Id:    newUserId,
		Name:  cmd.Name,
		Email: cmd.Email,
	})
	tokenPair, err := handler.jwtToken.GenerateTokenPair(newUserId.String(), "")
	if err != nil {
		return nil, generateTokenErr
	}
	response := responses.TokenResponse{
		UserId:       newUserId,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
	return &response, nil
}
