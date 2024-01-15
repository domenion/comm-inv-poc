package services

import (
	"comm-inv-poc/src/internal/core/models"
	"testing"
)

func TestImportInventories(t *testing.T) {
	svc := NewCMTService()
	tests := []struct {
		name string
		req  *models.ImportInventoriesRequest
	}{
		{
			name: "import",
			req: &models.ImportInventoriesRequest{
				Type: "inventory",
				Resources: []models.Resource{
					{
						Sku:             "80000084_3000010827",
						QuantityOnStock: 11,
					},
				},
			},
		},
	}
	token, err := svc.GetToken("")
	if err != nil {
		t.Error(err.Error())
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := svc.ImportInventories(token, tt.req)
			if err != nil {
				t.Error(err.Error())
				return
			}
			if resp == nil {
				t.Error("Response nil")
				return
			}
		})
	}
}
