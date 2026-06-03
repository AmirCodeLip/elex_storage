package client_repositories

import "net/http"

type FileStorageRepository interface {
	Upload(r *http.Request) (*http.Response, error)
}
