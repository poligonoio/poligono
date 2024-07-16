package services

type InfisicalService interface {
	GetSecret(key string) (string, error)
	CreateSecret(key string, secret string) error
	UpdateSecret(key string, secret string) error
	DeleteSecret(key string) error
}
