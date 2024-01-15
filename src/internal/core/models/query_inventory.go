package models

type QueryInventoryResponse struct {
	StatusCode int                    `json:"statusCode"`
	Message    string                 `json:"message"`
	Limit      int64                  `json:"limit"`
	Offset     int64                  `json:"offset"`
	Count      int64                  `json:"count"`
	Total      int64                  `json:"total"`
	Results    []QueryInventoryResult `json:"results"`
}

type QueryInventoryResult struct {
	ID                string `json:"id"`
	Version           int64  `json:"version"`
	Sku               string `json:"sku"`
	Key               string `json:"key"`
	QuantityOnStock   int64  `json:"quantityOnStock"`
	AvailableQuantity int64  `json:"availableQuantity"`
	CreatedAt         string `json:"createdAt"`
	LastModifiedAt    string `json:"lastModifiedAt"`
}
