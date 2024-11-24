package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/reversersed/zyfra/docs"
	"github.com/reversersed/zyfra/internal/config"
	"github.com/reversersed/zyfra/internal/handlers"
	"github.com/reversersed/zyfra/internal/service"
	"github.com/reversersed/zyfra/internal/storage"
	"github.com/reversersed/zyfra/pkg/mongo"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title SSO API
// @version 1.0

// @host localhost:9000
// @BasePath /api/

// @scheme http
// @accept json
type app struct {
	handler handlers.Handler
	cfg     *config.Config
}

func New() *app {
	flag.Parse()

	cfg, err := config.GetConfig("config/.env")
	if err != nil {
		log.Fatalf("couldn't open and read file: %v", err)
	}

	driver, err := mongo.NewClient(context.Background(), cfg.Database)
	if err != nil {
		log.Fatalf("couldn't create mongo connection: %v", err)
	}

	return &app{handler: handlers.New(service.New(storage.New(driver))), cfg: cfg}
}
func (a *app) Run() {
	engine := gin.Default()
	gin.SetMode(gin.DebugMode)

	a.handler.Register(engine)

	log.Printf("starting server listening to http://%s:%d...\n", a.cfg.Server.Host, a.cfg.Server.Port)
	engine.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := engine.Run(fmt.Sprintf("%s:%d", a.cfg.Server.Host, a.cfg.Server.Port)); err != nil {
		log.Fatalf("router error: %v", err)
	}
}
func (a *app) Close() {
	os.Exit(0)
}
