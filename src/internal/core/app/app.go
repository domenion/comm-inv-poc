package app

import (
	"comm-inv-poc/src/internal/core/models"
	"comm-inv-poc/src/internal/core/ports/repositories"
	"comm-inv-poc/src/internal/core/ports/services"
	"fmt"
	"log"
	"time"
)

type App struct {
	repo repositories.ProductRepository
	cmt  services.CMTService
}

func (a *App) ImportByItem() error {
	tst := time.Now()
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	qst := time.Now()
	productList, err := a.repo.GetProducts()
	if err != nil {
		log.Print(err.Error())
		return err
	}
	qet := time.Now()
	fmt.Printf("querying time: %f seconds\n", qet.Sub(qst).Seconds())

	authResp, err := a.cmt.Authenticate("manage_project")
	if err != nil {
		return err
	}

	auth := authResp.CreateToken()
	for _, p := range productList {
		st := time.Now()
		if ext, _ := a.cmt.GetInventoryByKey(auth, p.GetKey()); ext == nil {
			r, _ := a.cmt.CreateInventory(auth, p.ToCreateInventoryRequest())
			fmt.Println(r.Message)
		} else {
			r, _ := a.cmt.UpdateInventory(auth, p.ToUpdateInventoryRequest(ext.Version), ext.Id)
			fmt.Println(r.Message)
		}
		et := time.Now()
		fmt.Printf("importing time: %f seconds\n", et.Sub(st).Seconds())
	}
	tet := time.Now()
	fmt.Printf("total execution time: %f seconds\n", tet.Sub(tst).Seconds())
	fmt.Printf("total %d records", len(productList))
	return nil
}

func (a *App) ImportByContainer() error {
	startTime := time.Now()
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	sqt := time.Now()
	productList, err := a.repo.GetProducts()
	if err != nil {
		log.Print(err.Error())
		return err
	}

	eqt := time.Now()
	fmt.Print("query time: ", eqt.Sub(sqt).Seconds(), " seconds, ", len(productList), " records")

	containerAuthResp, err := a.cmt.Authenticate("manage_import_containers")
	if err != nil {
		return err
	}

	containerAuth := containerAuthResp.CreateToken()
	if !a.cmt.CheckImportContainer(containerAuth) {
		a.cmt.CreateImportContainer(containerAuth)
	}

	productsAuthResp, err := a.cmt.Authenticate("manage_products")
	if err != nil {
		return err
	}

	productsAuth := productsAuthResp.CreateToken()
	count := 0
	importRequest := models.NewImportInventoriesRequest()
	n := len(productList)
	for i := 0; i < n; i++ {
		prd := productList[i]
		resource := models.Resource{
			Key:             prd.GetKey(),
			Sku:             prd.ProductCode,
			QuantityOnStock: int64(prd.QTY),
			SupplyChannel: models.SupplyChannel{
				TypeID: "channel",
				Key:    prd.ShopCode,
			},
		}
		importRequest.Resources = append(importRequest.Resources, resource)
		reachedMax := i+1 == n
		count++
		if count == 20 || reachedMax {
			sit := time.Now()
			resp, err := a.cmt.ImportInventories(productsAuth, &importRequest)
			if err != nil {
				return err
			}

			eit := time.Now()
			fmt.Print("import time: ", eit.Sub(sit).Seconds(), " seconds")

			fmt.Print(resp.StatusCode, resp.Error, resp.Message)
			importRequest = models.NewImportInventoriesRequest()
			count = 0
		}
	}

	endTime := time.Now()
	fmt.Printf("total processing time: %f seconds", endTime.Sub(startTime).Seconds())
	fmt.Printf("total %d records", len(productList))
	return nil
}

func New(repo repositories.ProductRepository, cmt services.CMTService) *App {
	return &App{repo, cmt}
}
