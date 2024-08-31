package services

type EngineService interface {
	RemoveCatalog(catalogName string) error
	Query(query string, dest interface{}) error
	GetRawData(query string) ([]map[string]interface{}, error)
}
