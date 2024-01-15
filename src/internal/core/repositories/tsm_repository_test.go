package repositories

import "testing"

func TestGetProduct(t *testing.T) {
	repo := NewTSMRepository()
	tests := []struct {
		name string
	}{
		{name: "PRODUCT_CODE: 3000010827"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := repo.GetProduct("3000010827")
			if err != nil {
				t.Error(err.Error())
				return
			}
			if len(p) == 0 {
				t.Error("Product not found")
				return
			}
		})
	}
}
