package core_utils

import (
	"bytes"
	"elex_storage/file_storage/internal/domain/entities"
	pkg_entities "elex_storage/pkg/shared_kernel/entities"
	"elex_storage/pkg/shared_kernel/models"
	"path/filepath"
)

type PathUtil struct {
	config *models.ConfigEnv2
}

func NewPathUtil(config *models.ConfigEnv2) *PathUtil {
	return &PathUtil{config}
}

func (pathUtil *PathUtil) GetPath(fileEntity entities.FileEntity) (path string, fullPath string) {
	path = filepath.Join(
		pathUtil.config.DriveDisk,
		fileEntity.CreatedAt.Format("2006"),
		fileEntity.CreatedAt.Format("01"),
	)
	name := fileEntity.Id.String() + ".esx"
	fullPath = filepath.Join(path, name)
	return path, fullPath
}

func (pathUtil *PathUtil) GetContentType(data *[]byte) pkg_entities.FileContentType {
	d := *data
	size := len(d)
	if size < 20 {
		return pkg_entities.Unown
	}
	header := d[0:20]
	switch {
	case bytes.HasPrefix(header, []byte("%PDF-")):
		return pkg_entities.PDF
	case bytes.HasPrefix(header, []byte{0x89, 'P', 'N', 'G'}):
		return pkg_entities.PNG
	case bytes.HasPrefix(header, []byte{0xFF, 0xD8, 0xFF}):
		return pkg_entities.JPEG
	case len(header) >= 8 && string(header[4:8]) == "ftyp":
		return pkg_entities.MP4
	case bytes.HasPrefix(header, []byte{0x1A, 0x45, 0xDF, 0xA3}):
		return pkg_entities.MKV
	default:
		return pkg_entities.Unown
	}
}
