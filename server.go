package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
	docs "github.com/poligonoio/vega-core/docs"
	"github.com/poligonoio/vega-core/internal/controllers"
	"github.com/poligonoio/vega-core/internal/middlewares"
	"github.com/poligonoio/vega-core/internal/services"
	"github.com/poligonoio/vega-core/pkg/env"
	"github.com/poligonoio/vega-core/pkg/logger"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/api/option"
)

var server *http.Server
var version string

func init() {
	// version := os.Getenv("POLIGONO_VERSION")
	basePathStr := "/v1alpha1"

	port := os.Getenv("PORT")

	if port == "" {
		port = "8888"
	}

	logger.Info.Println("Initializing API server...")

	ctx := context.TODO()
	var err error

	// Connect to mongo
	logger.Info.Println("Connecting to MongoDB...")
	mongoconn := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	mongoClient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		logger.Error.Fatalf("Couldn't connect to MongoDB: %v\n", err)
	}

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		logger.Error.Fatalf("Couldn't connect to MongoDB: %v\n", err)
	}

	logger.Info.Println("MongoDB connection established successfully!")

	// Initialize validator
	validate := validator.New()

	// Initialize Gemini client
	logger.Info.Println("Creating client for Gemini...")
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		logger.Error.Fatalf("Couldn't create a client for gemini: %v\n", err)
	}

	model := client.GenerativeModel("gemini-1.5-pro")
	logger.Info.Println("Gemini client created successfully successfully!")

	// Initialize probe service and controller
	logger.Info.Println("Initializing Probe service and controller...")
	probeService := services.NewProbeService(ctx)
	probeController := controllers.NewProbeController(probeService)
	logger.Info.Println("Probe service and controller initialized successfully!")

	// Initialize Infisical service
	logger.Info.Println("Initializing Infisical...")
	infisicalService, err := services.NewInfisicalService(ctx, os.Getenv("INFISICAL_PROJECT_ID"), "/")
	if err != nil {
		logger.Error.Fatalf("Error connecting to infisical: %v\n", err)
	}

	logger.Info.Println("Infisical initialized successfully!")

	logger.Info.Println("Initializing Trino service")
	trinoService, err := services.NewTrinoService(ctx)
	if err != nil {
		logger.Error.Fatalf("Error connecting to trino: %v\n", err)
	}
	logger.Info.Println("Trino service initialized successfully!")

	//
	schemaCollection := mongoClient.Database("poligono").Collection("schemas")
	schemaService := services.NewSchemaService(ctx, schemaCollection)

	// Initialize Data source service and controller
	logger.Info.Println("Initializing Data source service and controller...")
	dataSourceCollection := mongoClient.Database("poligono").Collection("datasources")
	dataSourceService := services.NewDataSourceService(ctx, dataSourceCollection, infisicalService, trinoService, schemaService)
	dataSourceController := controllers.NewDataSourceController(dataSourceService, trinoService, schemaService, validate)
	logger.Info.Println("Data source service and controller Initialized successfully!")

	// Initialize core service and controller
	logger.Info.Println("Initializing Core service and controller...")
	coreService := services.NewCoreService(ctx, model)
	coreController := controllers.NewCoreController(coreService, dataSourceService, trinoService)
	logger.Info.Println("Core service and controller initialized successfully!")

	// Initialize Gin routes and and middlewares
	logger.Info.Println("Initializing Gin routes and middleware...")
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", port)
	docs.SwaggerInfo.BasePath = basePathStr

	probePath := router.Group("/probe")
	basePath := router.Group(basePathStr)

	basePath.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	enableAuth := env.GetBoolEnv("ENABLE_AUTHENTICATION")
	enableRBAC := env.GetBoolEnv("AUTHENTICATION_ENABLE_RBAC")
	authType := os.Getenv("AUTHENTICATION_TYPE")

	if enableAuth && authType == "BasicAuth" {
		basePath.Use(gin.BasicAuth(gin.Accounts{
			"poligono": "poligono",
		}))
	}

	if enableAuth && authType == "oauth2" {
		basePath.Use(middlewares.EnsureValidToken())
	}

	if enableRBAC && enableAuth && authType == "oauth2" {
		basePath.Use(middlewares.EnsureValidRole())
	}

	basePath.Use(middlewares.SetVarsToContext())

	logger.Info.Println("Gin routes and middleware Initialized successfully!")

	// Register Gin routes
	logger.Info.Println("Registering Gin routes...")
	coreController.RegisterCoreRoutes(basePath)
	probeController.RegisterProbeRoutes(probePath)
	dataSourceController.RegisterDataSourceRoutes(basePath)
	logger.Info.Println("Gin routes registered successfully!")

	// Initialize server
	server = &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        router,
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// @title           Poligono API
// @description     Democratizing data access through plain English.
// @termsOfService  https://swagger.io/terms/

// @contact.name   Poligono Support
// @contact.url    https://www.swagger.io/support
// @contact.email  dev@poligono.xyz

// @license.name  GNU Affero General Public License version 3
// @license.url   https://www.gnu.org/licenses/agpl-3.0.html

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Start server
	logger.Info.Printf("Starting server on port %s...\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		logger.Error.Fatalf("Server error: %v\n", err)
	}
}
