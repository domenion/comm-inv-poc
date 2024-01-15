package repositories

import (
	"comm-inv-poc/src/internal/configs"
	"comm-inv-poc/src/internal/entities"
	"context"
	"database/sql/driver"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
)

type TSMRepository interface {
	GetProduct(code string) ([]*entities.Product, error)
	GetProducts() ([]*entities.Product, error)
}

type tsmRepository struct{}

func NewTSMRepository() TSMRepository {
	return &tsmRepository{}
}

func (repo *tsmRepository) GetProduct(code string) ([]*entities.Product, error) {
	conn, err := repo.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ps, err := repo.GetProducts()
	if err != nil {
		return nil, err
	}

	pp := []*entities.Product{}
	for _, p := range ps {
		if p.ProductCode == code {
			pp = append(ps, p)
		}
	}

	return pp, nil
}

func (repo *tsmRepository) GetProducts() ([]*entities.Product, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	ps := []*entities.Product{}
	conn, err := repo.createConnection()
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	interval := time.Duration(cfg.TSM.Interval)
	names := []driver.NamedValue{
		{Name: "SHOP_CODE", Value: "", Ordinal: 1},
		{Name: "PRODUCT_CODE", Value: "", Ordinal: 2},
		{Name: "QTY", Value: 0, Ordinal: 3},
		{Name: "lastDate", Value: time.Now().Add(-time.Minute * interval), Ordinal: 4},
	}
	stmt := `
		SELECT SHOP_CODE, PRODUCT_CODE, QTY 
		FROM TBL_SHOP_PRODUCT_LOCATION
		WHERE LOCATION_CODE = 'MAIN_STOCK'
			AND UPDATE_DATE >= :lastDate
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ds := conn.QueryRowContext(ctx, stmt, names)
	defer ds.Close()

	for ds.Next_() {
		p := entities.Product{}
		err = ds.Scan(&p)
		if err != nil {
			continue
		}
		ps = append(ps, &p)
	}

	return ps, err
}

func (r *tsmRepository) createConnection() (*go_ora.Connection, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	connStr := go_ora.BuildUrl(cfg.TSM.Host, cfg.TSM.Port, cfg.TSM.Database, cfg.TSM.Username, cfg.TSM.Password, nil)
	conn, err := go_ora.NewConnection(connStr)
	if err != nil {
		return nil, err
	}

	err = conn.Open()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
