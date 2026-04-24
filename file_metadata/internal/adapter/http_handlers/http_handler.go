package http_handlers

import (
	"elex_storage/pkg/shared_kernel/utils"

	"elex_storage/file_metadata/internal/use_case/cqrs/command_handlers"

	"github.com/jmoiron/sqlx"

	"elex_storage/pkg/logger"
)

type HttpHandler struct {
	logger                 logger.Logger
	db                     *sqlx.DB
	httpErrorUtils         *utils.HttpErrorUtils
	createDirectoryHandler *command_handlers.CreateDirectoryHandler
}

func NewHttpHandler(logger logger.Logger, db *sqlx.DB, httpErrorUtils *utils.HttpErrorUtils, createDirectoryHandler *command_handlers.CreateDirectoryHandler) *HttpHandler {
	return &HttpHandler{logger: logger, db: db, httpErrorUtils: httpErrorUtils, createDirectoryHandler: createDirectoryHandler}
}
