package repositories

import "comm-inv-poc/src/internal/core/entities"

type ProductRepository interface {
	GetProducts() ([]*entities.Product, error)
}
