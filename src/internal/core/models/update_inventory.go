package models

type UpdateInventoryRequest struct {
	Version int64         `json:"version"`
	Actions []interface{} `json:"actions"`
}

type Quantity struct {
	Action   string `json:"action"`
	Quantity int64  `json:"quantity"`
}

type SetSupplyChannel struct {
	Action        string        `json:"action"`
	SupplyChannel SupplyChannel `json:"supplyChannel"`
}

type UpdateInventoryResponse struct {
	ID         string `json:"id"`
	Key        string `json:"key"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message,omitempty"`
}
