package adapters

import (
	"comm-inv-poc/src/internal/configs"
	"comm-inv-poc/src/internal/core/entities"
	"comm-inv-poc/src/internal/core/ports/repositories"
	"context"
	"database/sql/driver"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
)

type tsmAdapter struct{}

func (a *tsmAdapter) createConnection() (*go_ora.Connection, error) {
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

func (a *tsmAdapter) GetProducts() ([]*entities.Product, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	ps := []*entities.Product{}
	conn, err := a.createConnection()
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

func NewTSMAdapter() repositories.ProductRepository {
	return &tsmAdapter{}
}
