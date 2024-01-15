package manager

import (
	"comm-inv-poc/src/internal/core/models"
	"comm-inv-poc/src/internal/core/repositories"
	"comm-inv-poc/src/internal/core/services"
	"fmt"
	"log"
	"time"
)

func createToken(importAuthResp *models.GetTokenResponse) string {
	return fmt.Sprintf("%s %s", importAuthResp.TokenType, importAuthResp.AccessToken)
}

func Start() error {
	startTime := time.Now()
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	sqt := time.Now()
	repo := repositories.NewTSMRepository()
	productList, err := repo.GetProducts()
	if err != nil {
		log.Print(err.Error())
		return err
	}

	eqt := time.Now()
	fmt.Print("query time: ", eqt.Sub(sqt).Seconds(), " seconds, ", len(productList), " records")

	svc := services.NewCMTService()
	containerAuthResp, err := svc.Authenticate("manage_import_containers")
	if err != nil {
		return err
	}

	containerAuth := createToken(containerAuthResp)
	if !svc.CheckImportContainer(containerAuth) {
		svc.CreateImportContainer(containerAuth)
	}

	productsAuthResp, err := svc.Authenticate("manage_products")
	if err != nil {
		return err
	}

	productsAuth := createToken(productsAuthResp)
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
			resp, err := svc.ImportInventories(productsAuth, &importRequest)
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
	fmt.Print("total processing time: ", endTime.Sub(startTime).Seconds())
	return nil
}
