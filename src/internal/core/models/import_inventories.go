package models

type ImportInventoriesRequest struct {
	Type      string     `json:"type"`
	Resources []Resource `json:"resources"`
}

func NewImportInventoriesRequest() ImportInventoriesRequest {
	return ImportInventoriesRequest{
		Type:      "inventory",
		Resources: []Resource{},
	}
}

type Resource struct {
	Key             string        `json:"key"`
	Sku             string        `json:"sku"`
	QuantityOnStock int64         `json:"quantityOnStock"`
	SupplyChannel   SupplyChannel `json:"supplyChannel"`
}

type Custom struct {
	Type SupplyChannel `json:"type"`
}

type ImportInventoriesResponse struct {
	StatusCode      int         `json:"statusCode"`
	Error           string      `json:"error"`
	Message         string      `json:"message"`
	OperationStatus interface{} `json:"operationStatus"`
}

type CreateImportContainerRequest struct {
	Key string `json:"key"`
}

type CreateImportContainerResponse struct {
	Key string `json:"key"`
}

type GetContainerResponse struct {
	Key string `json:"key"`
}
