package models

type InfisicalToken struct {
	AccessToken       string `json:"accessToken"`
	ExpiresIn         int    `json:"expiresIn"`
	AccessTokenMaxTTL int    `json:"accessTokenMaxTTL"`
	TokenType         string `json:"tokenType"`
}

type InfisicalCreateSecretRequestBody struct {
	WorkspaceId           string `json:"workspaceId"`
	Environment           string `json:"environment"`
	SecretPath            string `json:"secretPath"`
	SecretValue           string `json:"secretValue"`
	SecretComment         string `json:"secretComment"`
	SkipMultilineEncoding bool   `json:"skipMultilineEncoding"`
	Type                  string `json:"type"`
}

type InfisicalUpdateSecretRequestBody struct {
	WorkspaceId           string `json:"workspaceId"`
	Environment           string `json:"environment"`
	SecretPath            string `json:"secretPath"`
	SecretValue           string `json:"secretValue"`
	SkipMultilineEncoding bool   `json:"skipMultilineEncoding"`
	Type                  string `json:"type"`
}

type InfisicalDeleteSecretRequestBody struct {
	WorkspaceId string `json:"workspaceId"`
	Environment string `json:"environment"`
	SecretPath  string `json:"secretPath"`
	Type        string `json:"type"`
}

type InfisicalSecret struct {
	Id            string `json:"id"`
	ID            string `json:"_id"`
	Workspace     string `json:"workspace"`
	Environment   string `json:"environment"`
	SecretKey     string `json:"secretKey"`
	SecretValue   string `json:"secretValue"`
	SecretComment string `json:"secretComment"`
	Version       int    `json:"version"`
	Type          string `json:"type"`
}

type InfisicalCreateSecretResponseBody struct {
	Secret InfisicalSecret `json:"secret"`
}

type InfisicalGetSecretResponseBody struct {
	Secret InfisicalSecret `json:"secret"`
}

type InfisicalUpdateSecretResponseBody struct {
	Secret InfisicalSecret `json:"secret"`
}

type InfisicalDeleteSecretResponseBody struct {
	Secret InfisicalSecret `json:"secret"`
}
