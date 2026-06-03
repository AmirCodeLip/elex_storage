package http_handlers

import (
	"elex_storage/file_storage/internal/use_case/cqrs/commands"
	"net/http"
)

func (h *HttpHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	/// Parse content file from form.
	read_err, fileData, name := h.readFile(r)
	if read_err != nil {
		h.logger.Error(read_err.Error())
		h.httpErrorUtils.BadRequest(w, read_err)
		return
	}
	cmd := commands.SaveFileCommand{Name: *name, Data: fileData}
	saveFileErr := h.saveFileHandler.Handle(cmd)
	if saveFileErr != nil {
		h.logger.Error(saveFileErr.Error())
		h.httpErrorUtils.BadRequest(w, saveFileErr)
	}
}
