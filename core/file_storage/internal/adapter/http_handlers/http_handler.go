package http_handlers

import (
	"elex_storage/file_storage/internal/use_case/cqrs/command_handlers"
	"elex_storage/pkg/logger"
	"elex_storage/pkg/shared_kernel/utils"
)

type HttpHandler struct {
	logger          logger.Logger
	saveFileHandler *command_handlers.SaveFileHandler
	httpErrorUtils  *utils.HttpErrorUtils
}

func NewHttpHandler(logger logger.Logger, saveFileHandler *command_handlers.SaveFileHandler, httpErrorUtils *utils.HttpErrorUtils) *HttpHandler {
	return &HttpHandler{logger: logger, saveFileHandler: saveFileHandler}
}
