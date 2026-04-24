package http_handlers

import (
	"elex_storage/pkg/shared_kernel/models"
	"errors"
	"io"
	"net/http"
)

var fileLimitErr = models.NewCommonError(errors.New("File too large"))

func (h *HttpHandler) readFile(r *http.Request) (rerr error, data *[]byte, name *string) {
	// Limit the size to 200MB
	err := r.ParseMultipartForm(1 << 30)
	if err != nil {
		return fileLimitErr, nil, nil
	}
	// Get file from body.
	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error(err.Error())
		return InvalidFileErr, nil, nil
	}
	defer file.Close()

	// Read file content as bytes
	file_data, read_err := io.ReadAll(file)
	if read_err != nil {
		h.logger.Error(err.Error())
		return InvalidFileErr, nil, nil
	}

	return nil, &file_data, &header.Filename
}
