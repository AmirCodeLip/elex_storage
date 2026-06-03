package http_handlers

import (
	"context"
	"elex_storage/pkg/shared_kernel/grpc_service"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (handler *HttpHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req grpc_service.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.httpErrorUtils.BadRequest(w, err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tokenResponse, err := (handler.UserServiceClient.Login(ctx, &req))
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			// Not a gRPC status error
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("unknown error: %v", err)))
		}
		code := st.Code()
		if code == codes.InvalidArgument || code == codes.Aborted {
			w.WriteHeader(http.StatusBadRequest)
			message := st.Message()
			w.Write([]byte(message))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("unknown error: %v", err)))
			return
		}
	}
	// handler.userServiceClient
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&tokenResponse); err != nil {
		handler.httpErrorUtils.BadRequest(w, err)
		return
	}
}
