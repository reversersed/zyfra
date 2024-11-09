package app

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/reversersed/zyfra/docs"
	"github.com/reversersed/zyfra/internal/config"
	"github.com/reversersed/zyfra/internal/handlers"
	"github.com/reversersed/zyfra/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/bcrypt"
)

var (
	fileConfig = flag.String("file", "", "Absolute path to config json file (with valid usernames)")
	userName   = flag.String("username", "admin", "Act like valid username if no config file provided")
	password   = flag.String("password", "admin", "Act like valid user password if no config file provided")
	port       = flag.Int("port", 80, "Port the server will be listening")
	help       = flag.Bool("help", false, "Prints flags manual")
)

const (
	COMMAND_LOGIN  = "login"
	ARGS_LOGIN     = "[username] [password]"
	COMMAND_QUIT   = "exit"
	COMMAND_AUTH   = "auth"
	ARGS_AUTH      = "[session key]"
	COMMAND_DELETE = "delete"
	ARGS_DELETE    = "[session key]"
)

// @title SSO API
// @version 1.0

// @host localhost:80
// @BasePath /api/

// @scheme http
// @accept json
type app struct {
	handler handlers.Handler
}

func New() *app {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return nil
	}

	cfg, err := config.ReadFromFile(*fileConfig)
	if err != nil && len(*fileConfig) > 0 {
		log.Fatalf("couldn't open and read file: %s (%v)", *fileConfig, err)
	} else if err != nil {
		log.Printf("config file not found, using %s for username and %s for password...", *userName, *password)
		pass, _ := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		cfg = map[string][]byte{*userName: pass}
	}

	return &app{handler: handlers.New(service.New(), cfg)}
}
func (a *app) Run() {
	engine := gin.Default()
	gin.SetMode(gin.DebugMode)

	a.handler.Register(engine)

	log.Printf("starting server listening to http://localhost:%d...", *port)
	engine.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Printf("swagger uri: http://localhost:%d/api/swagger/index.html", *port)
	if err := engine.Run(fmt.Sprintf("localhost:%d", *port)); err != nil {
		log.Fatalf("router error: %v", err)
	}
}
func (a *app) Close() {
	os.Exit(0)
}
