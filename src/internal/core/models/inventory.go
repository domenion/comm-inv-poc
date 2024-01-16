package models

type GetInventoryResponse struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Version int64  `json:"version"`
}
