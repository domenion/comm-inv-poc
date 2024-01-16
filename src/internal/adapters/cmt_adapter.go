package adapters

import (
	"comm-inv-poc/src/internal/configs"
	"comm-inv-poc/src/internal/core/models"
	"comm-inv-poc/src/internal/core/ports/services"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const CONTENT_TYPE_KEY = "Content-Type"
const CONTENT_TYPE_JSON = "application/json"
const IMPORT_CONTAINER_KEY = "import-inventories"

type cmtAdapter struct {
}

func httpRequest[T any](method string, auth string, url string, req *any) (code int, resp *T, err error) {
	var reqBodyRd *strings.Reader
	if req != nil {
		reqBodyRw, err := json.Marshal(req)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}

		reqBodyRd = strings.NewReader(string(reqBodyRw))
	} else {
		reqBodyRd = nil
	}

	var httpReq *http.Request
	if reqBodyRd != nil {
		httpReq, err = http.NewRequest(method, url, reqBodyRd)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
	} else {
		httpReq, err = http.NewRequest(method, url, nil)
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
	}

	httpReq.Header.Add(CONTENT_TYPE_KEY, CONTENT_TYPE_JSON)
	httpReq.Header.Add("Authorization", auth)
	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	defer httpResp.Body.Close()
	respBodyRw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return httpResp.StatusCode, nil, err
	}

	resp = new(T)
	err = json.Unmarshal(respBodyRw, resp)
	if err != nil {
		return httpResp.StatusCode, nil, err
	}

	return httpResp.StatusCode, resp, nil
}

func (a *cmtAdapter) createToken(clientId, secret string) string {
	data := fmt.Sprintf("%s:%s", clientId, secret)
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func (a *cmtAdapter) GetToken(target string) (string, error) {
	authResp, err := a.Authenticate(target)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", authResp.TokenType, authResp.AccessToken), nil
}

func (a *cmtAdapter) Authenticate(target string) (*models.GetTokenResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.AUTH
	service := "oauth/token"
	clientId := cfg.Commercetools.ClientID
	secret := cfg.Commercetools.ClientSecret
	scope := fmt.Sprintf("%s:%s", target, cfg.Commercetools.ProjectKey)
	token := a.createToken(clientId, secret)
	auth := fmt.Sprintf("Basic %s", token)
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", scope)

	uri, err := url.ParseRequestURI(server)
	if err != nil {
		return nil, err
	}

	uri.Path = service
	req, err := http.NewRequest(http.MethodPost, uri.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add(CONTENT_TYPE_KEY, "application/x-www-form-urlencoded")
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

func (*cmtAdapter) CheckImportContainer(auth string) bool {
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

func (*cmtAdapter) CreateImportContainer(auth string) error {
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

func (*cmtAdapter) ImportInventories(auth string, req *models.ImportInventoriesRequest) (*models.ImportInventoriesResponse, error) {
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

func (*cmtAdapter) CheckInventoryExist(auth string, key string) (*models.GetInventoryResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.API
	projectKey := cfg.Commercetools.ProjectKey
	service := "inventory"
	url := fmt.Sprintf(`%s/%s/%s/key=%s`, server, projectKey, service, key)
	httpReq, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", auth)
	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	var resp *models.GetInventoryResponse
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *cmtAdapter) GetInventoryByKey(auth string, key string) (*models.GetInventoryResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.API
	projectKey := cfg.Commercetools.ProjectKey
	service := "inventory"
	url := fmt.Sprintf("%s/%s/%s/key=%s", server, projectKey, service, key)
	code, resp, err := httpRequest[models.GetInventoryResponse](http.MethodGet, auth, url, nil)
	if err != nil {
		return nil, err
	}

	if code == http.StatusNotFound {
		return nil, nil
	}

	return resp, nil
}

func (a *cmtAdapter) CreateInventory(auth string, req *models.CreateInventoryRequest) (*models.CreateInventoryResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.API
	projectKey := cfg.Commercetools.ProjectKey
	service := "inventory"
	url := fmt.Sprintf("%s/%s/%s", server, projectKey, service)
	reqBodyRw, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqBodyRd := strings.NewReader(string(reqBodyRw))
	httpReq, err := http.NewRequest(http.MethodPost, url, reqBodyRd)
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

	resp := &models.CreateInventoryResponse{}
	err = json.Unmarshal(respBodyRw, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (*cmtAdapter) UpdateInventory(auth string, req *models.UpdateInventoryRequest, id string) (*models.UpdateInventoryResponse, error) {
	cfg, err := configs.GetConfig()
	if err != nil {
		return nil, err
	}

	server := cfg.Commercetools.API
	projectKey := cfg.Commercetools.ProjectKey
	service := "inventory"
	url := fmt.Sprintf("%s/%s/%s/%s", server, projectKey, service, id)
	reqBodyRw, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqBodyRd := strings.NewReader(string(reqBodyRw))
	httpReq, err := http.NewRequest(http.MethodPost, url, reqBodyRd)
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

	resp := &models.UpdateInventoryResponse{}
	err = json.Unmarshal(respBodyRw, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewCMTAdapter() services.CMTService {
	return &cmtAdapter{}
}
