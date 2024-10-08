package controllers

import (
	"fmt"
	"net/http"

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
	EngineService     services.EngineService
}

func NewCoreController(coreService services.CoreService, dataSourceService services.DataSourceService, engineService services.EngineService) CoreController {
	return CoreController{
		CoreService:       coreService,
		DataSourceService: dataSourceService,
		EngineService:     engineService,
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
// @Param Prompt body models.GenerateQueryBody true "Prompt Object"
// @Success 200 {object} models.GenerateQueryActivity
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /prompts/generate [post]
func (self *CoreController) GenerateQuery(c *gin.Context) {
	var generateQueryBody models.GenerateQueryBody

	_ownerId, _ := c.Get("owner_id")
	ownerId, ok := _ownerId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "internal_server_error",
			Description: "Failed to read user identity",
		})
	}

	_sub, _ := c.Get("sub")
	sub, ok := _sub.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "internal_server_error",
			Description: "Failed to read user identity",
		})
	}

	if err := c.ShouldBindJSON(&generateQueryBody); err != nil {
		logger.Error.Printf("[%s][%s] Failed to read request body: %v\n", ownerId, sub, err)
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Failed to read request body.",
		})
		return
	}

	// Get Data source info and secret
	ds, err := self.DataSourceService.GetByName(generateQueryBody.DataSourceName, ownerId, false)
	if err != nil {
		logger.Error.Printf("[%s][%s] Data source provided is invalid: %v\n", ownerId, sub, err)
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Data source provided is invalid.",
		})
		return
	}

	// Extract schemas info from Data source
	schemas, err := self.DataSourceService.GetDataSourceSchemas(ds.ID)
	if err != nil {
		logger.Error.Printf("[%s][%s] Error extracting metadata: %v\n", ownerId, sub, err)
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Error extracting data source schema",
		})
		return
	}

	schemasYaml, _ := yaml.Marshal(&schemas)
	var mergedPrompt string = fmt.Sprintf("I have a PostgreSQL Catalog in Trino named data_source_%s with the following database schema:\n\n%s\n\nGive me an SQL Trino Query that provides the following information: %s\n\nTo accomplish the task correctly please consider including schema on the query and return only a query without additional text.", ds.ID.Hex(), string(schemasYaml), generateQueryBody.Text)
	logger.Info.Println(mergedPrompt)

	// Generate query
	query, err := self.CoreService.PromptGemini(mergedPrompt)
	if err != nil {
		logger.Error.Printf("[%s][%s] Error processing prompt: %v\n", ownerId, sub, err)
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "internal_server_error",
			Description: "Error processing prompt",
		})
		return
	}

	var results []map[string]interface{}

	if generateQueryBody.Execute {
		// Get data from Data source using generated query
		results, err = self.EngineService.GetRawData(query)
		if err != nil {
			logger.Error.Printf("[%s][%s] Failed to get data from data source: %v\n", ownerId, sub, err)
			c.JSON(http.StatusInternalServerError, models.HTTPError{
				Error:       "bad_request",
				Description: "Failed to get data from data source.",
			})
			return
		}
	}

	// save activity
	activity := models.GenerateQueryActivity{
		Prompt:         generateQueryBody.Text,
		Query:          query,
		Data:           results,
		UserId:         sub,
		MergedPrompt:   mergedPrompt,
		DataSourceName: ds.Name,
		DataSourceId:   ds.ID,
		OrganizationId: ownerId,
	}

	c.JSON(http.StatusOK, activity)
}

// @BasePath /

// PingExample godoc
// @Summary Improve SQL Query
// @Schemes
// @Description Improve an SQL query based on a natural language prompt
// @Tags query
// @Accept json
// @Produce json
// @Param Prompt body models.ImproveQueryBody true "Prompt Object"
// @Success 200 {object} models.ImproveQueryActivity
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /prompts/improve [post]
func (self *CoreController) ImproveQuery(c *gin.Context) {
	var improveQueryBody models.ImproveQueryBody

	_ownerId, _ := c.Get("owner_id")
	ownerId, ok := _ownerId.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "internal_server_error",
			Description: "Failed to read user identity",
		})
	}

	_sub, _ := c.Get("sub")
	sub, ok := _sub.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "internal_server_error",
			Description: "Failed to read user identity",
		})
	}

	if err := c.ShouldBindJSON(&improveQueryBody); err != nil {
		logger.Error.Printf("[%s][%s] Failed to read request body: %v\n", ownerId, sub, err)
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Failed to read request body.",
		})
		return
	}

	// Get Data source info and secret
	ds, err := self.DataSourceService.GetByName(improveQueryBody.DataSourceName, ownerId, false)
	if err != nil {
		logger.Error.Printf("[%s][%s] Data source provided is invalid: %v\n", ownerId, sub, err)
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Data source provided is invalid.",
		})
		return
	}

	// Extract schemas info from Data source
	schemas, err := self.DataSourceService.GetDataSourceSchemas(ds.ID)
	if err != nil {
		logger.Error.Printf("[%s][%s] Error extracting metadata: %v\n", ownerId, sub, err)
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Error extracting data source schema",
		})
		return
	}

	schemasYaml, _ := yaml.Marshal(&schemas)
	var mergedPrompt string = fmt.Sprintf("I have a PostgreSQL Catalog in Trino named data_source_%s with the following database schema:\n\n%s\n\nEnhance this SQL Trino query for improved readability and performance: %s\n\nTo accomplish the task correctly please consider including schema on the query and return only a query without additional text.", ds.ID.Hex(), string(schemasYaml), improveQueryBody.Query)
	logger.Info.Println(mergedPrompt)

	// Generate query
	query, err := self.CoreService.PromptGemini(mergedPrompt)
	if err != nil {
		logger.Error.Printf("[%s][%s] Error processing prompt: %v\n", ownerId, sub, err)
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "internal_server_error",
			Description: "Error processing prompt",
		})
		return
	}

	var results []map[string]interface{}

	if improveQueryBody.Execute {
		// Get data from Data source using generated query
		results, err = self.EngineService.GetRawData(query)
		if err != nil {
			logger.Error.Printf("[%s][%s] Failed to get data from data source: %v\n", ownerId, sub, err)
			c.JSON(http.StatusInternalServerError, models.HTTPError{
				Error:       "bad_request",
				Description: "Failed to get data from data source.",
			})
			return
		}
	}

	// save activity
	activity := models.ImproveQueryActivity{
		OriginalQuery:  improveQueryBody.Query,
		ImprovedQuery:  query,
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
	coreRoute := rg.Group("prompts")
	coreRoute.POST("/generate", self.GenerateQuery)
	coreRoute.POST("/improve", self.ImproveQuery)
}
