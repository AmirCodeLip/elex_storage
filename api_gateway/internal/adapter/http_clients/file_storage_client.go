package http_clients

import (
	"elex_storage/api_gateway/internal/domain/client_repositories"
	"elex_storage/pkg/shared_kernel/models"
	"elex_storage/pkg/shared_kernel/utils"
	"net/http"
	"time"
)

type fileStorageClient struct {
	client *utils.HttpApiClient
}

func NewFileStorageClient(cfg *models.ConfigEnv) client_repositories.FileStorageRepository {
	client := fileStorageClient{
		utils.NewAPIClient(cfg.FileStorageHttpUrl.FullAddress, 2*time.Minute),
	}
	return &client
}

func (f *fileStorageClient) Upload(r *http.Request) (*http.Response, error) {
	return f.client.Do(http.MethodPost, "upload", r)
}
