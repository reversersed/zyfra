package app

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/reversersed/zyfra/internal/config"
	"github.com/reversersed/zyfra/internal/parser"
	"github.com/reversersed/zyfra/internal/reader"
	"github.com/reversersed/zyfra/internal/service"
)

var (
	fileConfig = flag.String("file", "", "Absolute path to config json file (with valid usernames)")
	userName   = flag.String("username", "admin", "Act like valid username if no config file provided")
	password   = flag.String("password", "admin", "Act like valid user password if no config file provided")
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

type readerService interface {
	WaitKey() string
}
type parserService interface {
	ParseCommand(string) (string, []string, error)
}
type sessionService interface {
	Close() error
	CreateSession() string
	CheckSession(string) error
	Delete(key string) error
}
type app struct {
	reader     readerService
	parser     parserService
	service    sessionService
	validNames map[string]string
}

func New() *app {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return nil
	}

	cfg, err := config.Read(*fileConfig)
	if err != nil && len(*fileConfig) > 0 {
		log.Fatalf("couldn't open and read file: %s (%v)", *fileConfig, err)
	} else if err != nil {
		fmt.Printf("config file not found, using %s for username and %s for password...", *userName, *password)
		cfg = map[string]string{*userName: *password}
	}

	return &app{validNames: cfg, parser: parser.New(), reader: reader.New(), service: service.New()}
}
func (a *app) Run() {
	for {
		fmt.Printf("\nCommands available:\n\t%s %s\n\t%s %s\n\t%s %s\n\t%s\n", COMMAND_LOGIN, ARGS_LOGIN, COMMAND_AUTH, ARGS_AUTH, COMMAND_DELETE, ARGS_DELETE, COMMAND_QUIT)
		cmd, args, err := a.parser.ParseCommand(a.reader.WaitKey())
		if err != nil {
			fmt.Printf("! Parse error: %v\n", err)
			continue
		}

		switch cmd {
		case COMMAND_QUIT:
			a.Close()
		case COMMAND_LOGIN:
			if len(args) != 2 {
				fmt.Printf("\n\n! Excepted 2 arguments, but got %d\n\n", len(args))
				continue
			}
			password, ok := a.validNames[args[0]]
			if !ok {
				fmt.Println("\n\n! User does not exist\n ")
				continue
			}
			if password != args[1] {
				fmt.Println("\n\n! Incorrect password\n ")
				continue
			}
			fmt.Printf("\nGenerated session key: %s\nSession will be expired in 1 minute\n", a.service.CreateSession())
			continue
		case COMMAND_AUTH:
			if len(args) != 1 {
				fmt.Printf("\n\n! Excepted 1 arguments, but got %d\n\n", len(args))
				continue
			}
			if err := a.service.CheckSession(args[0]); err == nil {
				fmt.Println("\n\nUser successfully authorized")
				continue
			} else {
				fmt.Printf("\n\n%v\n", err)
				continue
			}
		case COMMAND_DELETE:
			if len(args) != 1 {
				fmt.Printf("\n\n! Excepted 1 arguments, but got %d\n\n", len(args))
				continue
			}
			if err := a.service.Delete(args[0]); err == nil {
				fmt.Println("\n\nSession deleted successfully")
				continue
			} else {
				fmt.Printf("\n\n%v\n", err)
				continue
			}
		default:
			fmt.Printf("\n\n! Command %s does not exist\n\n", cmd)
			continue
		}
	}
}
func (a *app) Close() {
	a.service.Close()
	os.Exit(0)
}
