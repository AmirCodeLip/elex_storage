package http_handlers

import (
	"elex_storage/file_metadata/internal/use_case/cqrs/commands"
	"encoding/json"
	"net/http"
)

func (handler *HttpHandler) CreateDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	var createDirectoryCommand commands.CreateDirectoryCommand
	if err := json.NewDecoder(r.Body).Decode(&createDirectoryCommand); err != nil {
		handler.httpErrorUtils.BadRequest(w, err)
		return
	}
	err, _ := handler.createDirectoryHandler.Handle(createDirectoryCommand)
	if err != nil {
		handler.httpErrorUtils.BadRequest(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`ok`))
}
