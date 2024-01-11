package services

import (
	"comm-inv-poc/src/internal/configs"
	"comm-inv-poc/src/internal/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CMTService interface {
	CheckImportContainer(auth string) bool
	CreateImportContainer(auth string) error
	QueryInventory(auth string, sku string) (*models.QueryInventoryResponse, error)
	CheckInventoryExist(auth string, key string) (int, error)
	CreateInventory(auth string, req *models.CreateInventoryRequest) (*models.CreateInventoryResponse, error)
	UpdateInventory(auth string, req *models.UpdateInventoryRequest, id string) (*models.UpdateInventoryResponse, error)
	ImportInventories(auth string, req *models.ImportInventoriesRequest) (*models.ImportInventoriesResponse, error)
	Authenticate(target string) (*models.GetTokenResponse, error)
	GetToken(target string) (string, error)
}

type cmtService struct{}

func NewCMTService() CMTService {
	return &cmtService{}
}

const CONTENT_TYPE_KEY = "Content-Type"
const CONTENT_TYPE_JSON = "application/json"
const IMPORT_CONTAINER_KEY = "import-inventories"

func (s *cmtService) CheckImportContainer(auth string) bool {
	cfg, err := configs.GetConfig()
	if err != nil {
		return false
	}

	server := cfg.Commercetools.Import
	projectKey := cfg.Commercetools.ProjectKey
	service := "import-containers"
	url := fmt.Sprintf(`%s/%s/%s/%s`, server, projectKey, service, IMPORT_CONTAINER_KEY)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	req.Header.Add(CONTENT_TYPE_KEY, CONTENT_TYPE_JSON)
	req.Header.Add("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode != http.StatusNotFound
}

func (s *cmtService) CreateImportContainer(auth string) error {
	cfg, err := configs.GetConfig()
	if err != nil {
		return err
	}

	server := cfg.Commercetools.Import
	projectKey := cfg.Commercetools.ProjectKey
	service := "import-containers"
	url := fmt.Sprintf(`%s/%s/%s`, server, projectKey, service)
	reqBody := &models.CreateImportContainerRequest{Key: IMPORT_CONTAINER_KEY}
	reqBodyRw, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	reqBodyRd := strings.NewReader(string(reqBodyRw))
	req, err := http.NewRequest(http.MethodPost, url, reqBodyRd)
	if err != nil {
		return err
	}

	req.Header.Add(CONTENT_TYPE_KEY, CONTENT_TYPE_JSON)
	req.Header.Add("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	respBodyRw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respBody := &models.CreateImportContainerResponse{}
	err = json.Unmarshal(respBodyRw, respBody)
	if err != nil {
		return err
	}

	return nil
}

func (s *cmtService) GetProduct(auth string, key string) (*models.GetProductResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.API
	projectKey := cfg.Commercetools.ProjectKey
	service := "products"
	url := fmt.Sprintf(`%s/%s/%s/%s`, server, projectKey, service, key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(CONTENT_TYPE_KEY, CONTENT_TYPE_JSON)
	req.Header.Add("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBodyRw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Print(string(respBodyRw))
	respBody := &models.GetProductResponse{}
	err = json.Unmarshal(respBodyRw, respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (s *cmtService) QueryInventory(auth string, sku string) (*models.QueryInventoryResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.API
	projectKey := cfg.Commercetools.ProjectKey
	service := "inventory"
	url := fmt.Sprintf(`%s/%s/%s`, server, projectKey, service)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(CONTENT_TYPE_KEY, CONTENT_TYPE_JSON)
	req.Header.Add("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBodyRw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respBody := &models.QueryInventoryResponse{}
	err = json.Unmarshal(respBodyRw, respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (s *cmtService) CheckInventoryExist(auth string, key string) (int, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	server := cfg.Commercetools.API
	projectKey := cfg.Commercetools.ProjectKey
	service := "inventory"
	url := fmt.Sprintf(`%s/%s/%s/key=%s`, server, projectKey, service, key)
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req.Header.Add("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	code := resp.StatusCode
	return code, nil
}

func (s *cmtService) CreateInventory(auth string, reqBody *models.CreateInventoryRequest) (*models.CreateInventoryResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.API
	projectKey := cfg.Commercetools.ProjectKey
	service := "inventory"
	url := fmt.Sprintf("%s/%s/%s", server, projectKey, service)
	reqBodyRw, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	reqBodyRd := strings.NewReader(string(reqBodyRw))
	req, err := http.NewRequest(http.MethodPost, url, reqBodyRd)
	if err != nil {
		return nil, err
	}

	req.Header.Add(CONTENT_TYPE_KEY, CONTENT_TYPE_JSON)
	req.Header.Add("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBodyRw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respBody := &models.CreateInventoryResponse{}
	err = json.Unmarshal(respBodyRw, respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (s *cmtService) UpdateInventory(auth string, reqBody *models.UpdateInventoryRequest, id string) (*models.UpdateInventoryResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.API
	projectKey := cfg.Commercetools.ProjectKey
	service := "inventory"
	url := fmt.Sprintf("%s/%s/%s/%s", server, projectKey, service, id)
	reqBodyRw, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	reqBodyRd := strings.NewReader(string(reqBodyRw))
	req, err := http.NewRequest(http.MethodPost, url, reqBodyRd)
	if err != nil {
		return nil, err
	}

	req.Header.Add(CONTENT_TYPE_KEY, CONTENT_TYPE_JSON)
	req.Header.Add("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBodyRw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respBody := &models.UpdateInventoryResponse{}
	err = json.Unmarshal(respBodyRw, respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (s *cmtService) ImportInventories(auth string, req *models.ImportInventoriesRequest) (*models.ImportInventoriesResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.Import
	projectKey := cfg.Commercetools.ProjectKey
	service := "inventories/import-containers"
	url := fmt.Sprintf("%s/%s/%s/%s", server, projectKey, service, IMPORT_CONTAINER_KEY)

	reqBodyRw, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqBodyRd := strings.NewReader(string(reqBodyRw))

	httpReq, err := http.NewRequest(http.MethodPost, url, reqBodyRd)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	httpReq.Header.Add(CONTENT_TYPE_KEY, CONTENT_TYPE_JSON)
	httpReq.Header.Add("Authorization", auth)
	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	respBodyRw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	resp := &models.ImportInventoriesResponse{}
	err = json.Unmarshal(respBodyRw, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *cmtService) Authenticate(target string) (*models.GetTokenResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.AUTH
	service := "oauth/token"
	clientId := cfg.Commercetools.ClientID
	secret := cfg.Commercetools.ClientSecret
	grantType := "client_credentials"
	scope := fmt.Sprintf("%s:%s", target, cfg.Commercetools.ProjectKey)
	url := fmt.Sprintf("%s/%s?grant_type=%s&scope=%s", server, service, grantType, scope)
	token := s.createToken(clientId, secret)
	auth := fmt.Sprintf("Basic %s", token)
	reqBody := &models.GetTokenRequest{}
	reqBodyRw, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	reqBodyRd := strings.NewReader(string(reqBodyRw))
	req, err := http.NewRequest(http.MethodPost, url, reqBodyRd)
	if err != nil {
		return nil, err
	}

	req.Header.Add(CONTENT_TYPE_KEY, CONTENT_TYPE_JSON)
	req.Header.Add("Authorization", auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBodyRw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respBody := &models.GetTokenResponse{}
	err = json.Unmarshal(respBodyRw, respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (s *cmtService) GetToken(target string) (string, error) {
	authResp, err := s.Authenticate(target)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", authResp.TokenType, authResp.AccessToken), nil
}

func (s *cmtService) createToken(clientId, secret string) string {
	data := fmt.Sprintf("%s:%s", clientId, secret)
	return base64.StdEncoding.EncodeToString([]byte(data))
}
