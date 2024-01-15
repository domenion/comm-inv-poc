package models

type CreateInventoryRequest struct {
	Sku               string             `json:"sku"`
	Key               string             `json:"key"`
	SupplyChannel     ResourceIdentifier `json:"supplyChannel"`
	QuantityOnStock   int64              `json:"quantityOnStock"`
	AvailableQuantity int64              `json:"availableQuantity"`
}

type ResourceIdentifier struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type CreateInventoryResponse struct {
	ID         string `json:"id"`
	Key        string `json:"key"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message,omitempty"`
}
