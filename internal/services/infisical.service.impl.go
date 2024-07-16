package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/poligonoio/vega-core/internal/models"
	"github.com/poligonoio/vega-core/pkg/logger"
)

var ErrSecretAlreadyExist = fmt.Errorf("Secret already exist")

type InfisicalErrorResponse struct {
	Error      string
	Message    string
	StatusCode int
}

type InfisicalServiceImpl struct {
	ctx        context.Context
	token      models.InfisicalToken
	projectId  string
	secretPath string
}

func NewInfisicalService(ctx context.Context, projectId string, secretPath string) (InfisicalService, error) {
	data := url.Values{}
	data.Set("clientId", os.Getenv("INFISICAL_CLIENT_ID"))
	data.Set("clientSecret", os.Getenv("INFISICAL_CLIENT_SECRET"))

	req, err := http.NewRequest("POST", "https://app.infisical.com/api/v1/auth/universal-auth/login", strings.NewReader(data.Encode()))
	if err != nil {
		logger.Error.Println(fmt.Printf("Failed to get token from Infisical: %v\n", err))
		return nil, err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error.Println(fmt.Printf("Failed to get token from Infisical: %v\n", err))
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Couldn't connect to Infisical: %s", res.Status)
	}

	token := models.InfisicalToken{}
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		logger.Error.Println(fmt.Printf("Failed to get decode token from Infisical: %v\n", err))
		return nil, err
	}

	return &InfisicalServiceImpl{
		ctx:        ctx,
		token:      token,
		projectId:  projectId,
		secretPath: secretPath}, nil
}

func (self *InfisicalServiceImpl) GetSecret(key string) (string, error) {
	url := fmt.Sprintf("https://app.infisical.com/api/v3/secrets/raw/%s?workspaceId=%s&environment=general", key, self.projectId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", self.token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Couldn't connect to Infisical: %s", res.Status)
	}

	getSecretResponseBody := models.InfisicalGetSecretResponseBody{}
	if err = json.NewDecoder(res.Body).Decode(&getSecretResponseBody); err != nil {
		return "", err
	}

	return getSecretResponseBody.Secret.SecretValue, nil
}

func (self *InfisicalServiceImpl) CreateSecret(key string, secret string) error {
	logger.Info.Println(self.token)
	url := fmt.Sprintf("https://app.infisical.com/api/v3/secrets/raw/%s", key)

	requestBody := models.InfisicalCreateSecretRequestBody{
		WorkspaceId:           self.projectId,
		SecretPath:            self.secretPath,
		SecretValue:           secret,
		Environment:           "general",
		SecretComment:         "",
		SkipMultilineEncoding: true,
		Type:                  "shared",
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", self.token.AccessToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		errRes := InfisicalErrorResponse{}

		if err = json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return err
		}

		if errRes.StatusCode == 400 && errRes.Message == "Secret already exist" {
			return ErrSecretAlreadyExist
		}

		logger.Error.Println(errRes)
		return fmt.Errorf("Failed to connect to Infisical: %s", res.Status)
	}

	createSecretResponseBody := models.InfisicalCreateSecretResponseBody{}
	if err = json.NewDecoder(res.Body).Decode(&createSecretResponseBody); err != nil {
		return err
	}

	fmt.Println(createSecretResponseBody.Secret)

	return nil
}

func (self *InfisicalServiceImpl) UpdateSecret(key string, secret string) error {
	url := fmt.Sprintf("https://app.infisical.com/api/v3/secrets/raw/%s", key)

	requestBody := models.InfisicalUpdateSecretRequestBody{
		WorkspaceId:           self.projectId,
		SecretPath:            self.secretPath,
		SecretValue:           secret,
		Environment:           "general",
		SkipMultilineEncoding: true,
		Type:                  "shared",
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", self.token.AccessToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Couldn't connect to Infisical: %s", res.Status)
	}

	updateSecretResponseBody := models.InfisicalUpdateSecretResponseBody{}
	if err = json.NewDecoder(res.Body).Decode(&updateSecretResponseBody); err != nil {
		return err
	}

	fmt.Println(updateSecretResponseBody.Secret)

	return nil
}

func (self *InfisicalServiceImpl) DeleteSecret(key string) error {
	url := fmt.Sprintf("https://app.infisical.com/api/v3/secrets/raw/%s", key)

	requestBody := models.InfisicalDeleteSecretRequestBody{
		WorkspaceId: self.projectId,
		SecretPath:  self.secretPath,
		Environment: "general",
		Type:        "shared",
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", self.token.AccessToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Couldn't to connect to infisical: %s", res.Status)
	}

	deleteSecretResponseBody := models.InfisicalDeleteSecretResponseBody{}
	if err = json.NewDecoder(res.Body).Decode(&deleteSecretResponseBody); err != nil {
		return err
	}

	fmt.Println(deleteSecretResponseBody.Secret)

	return nil
}
