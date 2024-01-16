package services

import "comm-inv-poc/src/internal/core/models"

type CMTService interface {
	GetToken(target string) (string, error)
	Authenticate(target string) (*models.GetTokenResponse, error)
	CheckImportContainer(auth string) bool
	CreateImportContainer(auth string) error
	ImportInventories(auth string, req *models.ImportInventoriesRequest) (*models.ImportInventoriesResponse, error)
	CheckInventoryExist(auth string, key string) (*models.GetInventoryResponse, error)
	GetInventoryByKey(auth string, key string) (*models.GetInventoryResponse, error)
	CreateInventory(auth string, req *models.CreateInventoryRequest) (*models.CreateInventoryResponse, error)
	UpdateInventory(auth string, req *models.UpdateInventoryRequest, id string) (*models.UpdateInventoryResponse, error)
}
