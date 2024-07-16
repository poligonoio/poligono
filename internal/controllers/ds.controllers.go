package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	validatorv10 "github.com/go-playground/validator/v10"
	"github.com/poligonoio/vega-core/internal/models"
	"github.com/poligonoio/vega-core/internal/services"
	"github.com/poligonoio/vega-core/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataSourceController struct {
	DataSourceService services.DataSourceService
	TrinoService      services.TrinoService
	SchemaService     services.SchemaService
	validate          *validatorv10.Validate
}

func NewDataSourceController(dataSourceService services.DataSourceService, trinoService services.TrinoService, schemaService services.SchemaService, validate *validatorv10.Validate) DataSourceController {
	return DataSourceController{
		DataSourceService: dataSourceService,
		TrinoService:      trinoService,
		SchemaService:     schemaService,
		validate:          validate,
	}
}

// @BasePath /

// PingExample godoc
// @Summary Retrieve Data Source by Name
// @Schemes
// @Description Retrieve data source configuration with the specified name
// @Tags data_source
// @Accept json
// @Produce json
// @Param name path string true "Data Source Name" string
// @Success 200 {object} models.DataSource
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Router /datasources/{name} [get]
func (self *DataSourceController) GetDataSourceByName(c *gin.Context) {
	name := c.Param("name")

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

	dataSource, err := self.DataSourceService.GetByName(name, ownerId)
	if err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] failed to get data source by name: %v\n", ownerId, sub, err))

		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, models.HTTPError{
				Error:       "not_found",
				Description: fmt.Sprintf("Data source '%s' not found", name),
			})
			return
		}

		c.JSON(http.StatusNotFound, models.HTTPError{
			Error:       "not_found",
			Description: "Failed to get data source by name",
		})
		return
	}

	dataSource.Secret = ""
	c.JSON(http.StatusOK, dataSource)
}

// @BasePath /

// PingExample godoc
// @Summary List Data Sources
// @Schemes
// @Description Retrieve all data sources associated with the specified organization
// @Tags data_source
// @Accept json
// @Produce json
// @Success 200 {array} models.DataSource
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Router /datasources/all [get]
func (self *DataSourceController) GetAllDataSource(c *gin.Context) {
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

	dataSources, err := self.DataSourceService.GetAll(ownerId)
	if err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] Couldn't find data: %v\n", ownerId, sub, err))
		c.JSON(http.StatusNotFound, models.HTTPError{
			Error:       "not_found",
			Description: "No data source found for your organization",
		})
		return
	}

	c.JSON(http.StatusOK, dataSources)
}

// @BasePath /

// PingExample godoc
// @Summary Add Data Source
// @Schemes
// @Description Create a new data source configuration
// @Tags data_source
// @Accept json
// @Produce json
// @Param Data_Source body models.DataSource true "Data Source"
// @Success 200 {object} models.HTTPSuccess
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /datasources [post]
func (self *DataSourceController) CreateDataSource(c *gin.Context) {
	var newDataSource models.DataSource

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

	if err := c.ShouldBindJSON(&newDataSource); err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] Failed to read request body: %v\n", ownerId, sub, err))
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Failed to read request body",
		})
		return
	}

	newDataSource.CreatedBy = sub
	newDataSource.OrganizationId = ownerId
	// newDataSource.UserInfo.Name = userData.Name
	// newDataSource.UserInfo.Picture = userData.Picture
	newDataSource.CreatedAt = time.Now()
	newDataSource.UpdatedAt = time.Now()

	if err := self.validate.Struct(newDataSource); err != nil {
		logger.Info.Println(newDataSource)
		validationErr := err.(validatorv10.ValidationErrors)
		logger.Error.Println(fmt.Printf("[%s][%s] One or more data source fields are invalid: %s\n", ownerId, sub, validationErr))
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "One or more data source fields are invalid.",
		})
		return
	}

	if err := self.DataSourceService.Create(newDataSource); err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] Failed to create data source: %v\n", ownerId, sub, err))
		c.JSON(http.StatusInternalServerError, models.HTTPError{
			Error:       "internal_server_error",
			Description: "Failed to create data source",
		})
		return
	}

	c.JSON(http.StatusOK, models.HTTPSuccess{Message: "success"})
}

// @BasePath /

// PingExample godoc
// @Summary Modify Data Source
// @Schemes
// @Description Update the configuration of the specified data source
// @Tags data_source
// @Accept json
// @Produce json
// @Param name path string true "Data Source Name" string
// @Param Data_Source body models.UpdateRequestDataSourceBody true "Data Source"
// @Success 200 {object} models.HTTPSuccess
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Router /datasources/{name} [put]
func (self *DataSourceController) UpdateDataSourceByName(c *gin.Context) {
	var updateDataSource models.UpdateRequestDataSourceBody

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

	name := c.Param("name")

	if err := c.ShouldBindJSON(&updateDataSource); err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] Failed to read request body: %v\n", ownerId, sub, err))
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Failed to read request body",
		})
		return
	}

	updateDataSource.UpdatedAt = time.Now()

	err := self.DataSourceService.Update(name, ownerId, updateDataSource)
	if err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] failed to update data source: %v\n", ownerId, sub, err))
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Failed to update data source",
		})
		return
	}

	c.JSON(http.StatusOK, models.HTTPSuccess{Message: "success"})
}

// @BasePath /

// PingExample godoc
// @Summary Remove Data Source
// @Schemes
// @Description Permanently deletes the specified data source
// @Tags data_source
// @Accept json
// @Produce json
// @Param name path string true "Data Source Name" string
// @Success 200 {object} models.HTTPSuccess
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /datasources/{name} [delete]
func (self *DataSourceController) DeleteDataSourceByName(c *gin.Context) {
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

	name := c.Param("name")

	err := self.DataSourceService.Delete(name, ownerId)
	if err != nil {
		logger.Error.Println(fmt.Printf("[%s][%s] Failed to delete data source: %v\n", ownerId, sub, err))
		c.JSON(http.StatusBadRequest, models.HTTPError{
			Error:       "bad_request",
			Description: "Failed to delete data source",
		})
		return
	}

	c.JSON(http.StatusOK, models.HTTPSuccess{Message: "success"})
}

func (self *DataSourceController) RegisterDataSourceRoutes(rg *gin.RouterGroup) {
	dataSourceRoute := rg.Group("datasources")
	dataSourceRoute.GET("/:name", self.GetDataSourceByName)
	dataSourceRoute.GET("/all", self.GetAllDataSource)
	dataSourceRoute.POST("", self.CreateDataSource)
	dataSourceRoute.PUT("/:name", self.UpdateDataSourceByName)
	dataSourceRoute.DELETE("/:name", self.DeleteDataSourceByName)
}
