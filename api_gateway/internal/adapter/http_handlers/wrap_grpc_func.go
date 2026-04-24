package http_handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WrapGrpcFunc[Req any, Res any](
	handler *HttpHandler,
	fn func(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Res, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handler.httpErrorUtils.BadRequest(w, errors.New("Can't parse content"))
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		response, err := (fn(ctx, &req))
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				// Not a gRPC status error
				handler.httpErrorUtils.InternalServerError(w)
			}
			code := st.Code()
			if code == codes.InvalidArgument || code == codes.Aborted {
				w.WriteHeader(http.StatusBadRequest)
				message := st.Message()
				w.Write([]byte(message))
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				handler.logger.Error(err.Error())
				handler.httpErrorUtils.InternalServerError(w)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			handler.httpErrorUtils.InternalServerError(w)
			return
		}
	}
}
