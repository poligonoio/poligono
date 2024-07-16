package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/poligonoio/vega-core/internal/models"
	"github.com/poligonoio/vega-core/internal/services"
	"github.com/poligonoio/vega-core/pkg/logger"
	"gopkg.in/yaml.v3"
)

type CoreController struct {
	CoreService       services.CoreService
	DataSourceService services.DataSourceService
	TrinoService      services.TrinoService
}

func NewCoreController(coreService services.CoreService, dataSourceService services.DataSourceService, trinoService services.TrinoService) CoreController {
	return CoreController{
		CoreService:       coreService,
		DataSourceService: dataSourceService,
		TrinoService:      trinoService,
	}
}

// @BasePath /

// PingExample godoc
// @Summary Generate SQL Query
// @Schemes
// @Description Create an SQL query based on a natural language prompt
// @Tags query
// @Accept json
// @Produce json
// @Param Prompt body models.Prompt true "Prompt Object"
// @Success 200 {object} models.Activity
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /prompt [post]
func (self *CoreController) GeminiPrompt(c *gin.Context) {
	var prompt models.Prompt
	var queryResult models.GeminiQueryResult

	_ownerId, _ := c.Get("owner_id")
	ownerId, ok := _ownerId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "internal_server_error",
			Description: "Failed to read user identity",
		})
	}

	_sub, _ := c.Get("owner_id")
	sub, ok := _sub.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "internal_server_error",
			Description: "Failed to read user identity",
		})
	}

	if err := c.ShouldBindJSON(&prompt); err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] Failed to read request body: %v\n", ownerId, sub, err))
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Failed to read request body.",
		})
		return
	}

	// Get Data source info and secret
	ds, err := self.DataSourceService.GetByName(prompt.DataSourceName, ownerId)
	if err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] Data source provided is invalid: %v\n", ownerId, sub, err))
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Data source provided is invalid.",
		})
		return
	}

	// Extract schemas info from Data source
	schemas, err := self.DataSourceService.GetDataSourceSchemas(ds.Name, ds.OrganizationId)
	logger.Info.Println(schemas)
	if err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] Error extracting metadata: %v\n", ownerId, sub, err))
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Error extracting data source schema",
		})
		return
	}

	schemasYaml, _ := yaml.Marshal(&schemas)
	catalogName := self.TrinoService.GetCatalogName(ds.Name, ds.OrganizationId)
	var mergedPrompt string = fmt.Sprintf("I have a PostgreSQL Catalog in Trino named %s with the following database schema:\n\n%s\n\nGive me an SQL Trino Query that provides the following information: %s\n\nTo accomplish the task correctly please consider including schema on the query and return only a query without additional text.", catalogName, string(schemasYaml), prompt.Text)
	logger.Info.Println(mergedPrompt)

	// Generate query
	queryResult, err = self.CoreService.GeminiPrompt(mergedPrompt)
	if err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] Error processing prompt: %v\n", ownerId, sub, err))
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "internal_server_error",
			Description: "Error processing prompt",
		})
		return
	}

	query := strings.ReplaceAll(queryResult.QueryMarkdown, "```", "")
	query = strings.ReplaceAll(query, "sql", "")

	// Get data from Data source using generated query
	results, err := self.TrinoService.GetRawData(query)
	if err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] Failed to get data from data source:: %v\n", ownerId, sub, err))
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "bad_request",
			Description: "Failed to get data from data source.",
		})
		return
	}

	// save activity
	activity := models.Activity{
		Prompt:         prompt.Text,
		Query:          queryResult.QueryMarkdown,
		Data:           results,
		UserId:         sub,
		MergedPrompt:   mergedPrompt,
		DataSourceName: ds.Name,
		DataSourceId:   ds.ID,
		OrganizationId: ownerId,
	}

	c.JSON(http.StatusOK, activity)
}

func (self *CoreController) RegisterCoreRoutes(rg *gin.RouterGroup) {
	coreRoute := rg.Group("prompt")
	coreRoute.POST("", self.GeminiPrompt)
}
