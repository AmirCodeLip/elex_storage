package command_handlers

import (
	"elex_storage/file_metadata/internal/domain/repositories"
	"elex_storage/file_metadata/internal/use_case/dtos"
	"elex_storage/pkg/logger"

	"github.com/jinzhu/copier"
)

type GetDirectoriesHandler struct {
	logger                      logger.Logger
	directoryMetadataRepository repositories.DirectoryMetadataRepository
}

func NewGetDirectoriesHandler(logger logger.Logger,
	directoryMetadataRepository repositories.DirectoryMetadataRepository) *GetDirectoriesHandler {
	return &GetDirectoriesHandler{logger, directoryMetadataRepository}
}

func (u *GetDirectoriesHandler) Handle() ([]dtos.DirectoryDto, error) {
	var result []dtos.DirectoryDto
	dirs, err := u.directoryMetadataRepository.GetDirectories()
	if err != nil {
		return nil, err
	}
	copier.Copy(&result, &dirs)
	return result, nil
}
