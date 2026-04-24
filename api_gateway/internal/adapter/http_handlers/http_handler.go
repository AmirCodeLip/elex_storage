package http_handlers

import (
	"elex_storage/pkg/shared_kernel/grpc_service"
	"elex_storage/pkg/shared_kernel/utils"

	"elex_storage/pkg/logger"
)

type HttpHandler struct {
	logger            logger.Logger
	httpErrorUtils    *utils.HttpErrorUtils
	UserServiceClient grpc_service.UserServiceClient
}

func NewHttpHandler(logger logger.Logger, httpErrorUtils *utils.HttpErrorUtils, userServiceClient grpc_service.UserServiceClient) *HttpHandler {
	return &HttpHandler{logger: logger, httpErrorUtils: httpErrorUtils, UserServiceClient: userServiceClient}
}
