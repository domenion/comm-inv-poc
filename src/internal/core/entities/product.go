package entities

import (
	"comm-inv-poc/src/internal/core/models"
	"fmt"
)

type Product struct {
	ShopCode    string `db:"SHOP_CODE,varchar,50"`
	ProductCode string `db:"PRODUCT_CODE,varchar,50"`
	QTY         int    `db:"QTY,number"`
}

func (t *Product) ToCreateInventoryRequest() *models.CreateInventoryRequest {
	return &models.CreateInventoryRequest{
		Sku: t.ProductCode,
		Key: t.GetKey(),
		SupplyChannel: models.ResourceIdentifier{
			Key:  t.ShopCode,
			Type: "channel",
		},
		QuantityOnStock:   int64(t.QTY),
		AvailableQuantity: int64(t.QTY),
	}
}

func (t *Product) ToUpdateInventoryRequest(ver int64) *models.UpdateInventoryRequest {
	changeQ := models.Quantity{
		Action:   "changeQuantity",
		Quantity: int64(t.QTY),
	}
	return &models.UpdateInventoryRequest{
		Version: ver,
		Actions: []interface{}{changeQ},
	}
}

func (t *Product) GetKey() string {
	return fmt.Sprintf("%s_%s", t.ShopCode, t.ProductCode)
}
