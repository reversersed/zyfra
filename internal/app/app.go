package app

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/reversersed/zyfra/internal/config"
	"github.com/reversersed/zyfra/internal/parser"
	"github.com/reversersed/zyfra/internal/reader"
	"github.com/reversersed/zyfra/internal/service"
	"golang.org/x/crypto/bcrypt"
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
	WaitForInput(io.Reader) string
}
type parserService interface {
	ParseCommand(string) (string, []string, error)
}
type sessionService interface {
	CreateSession() string
	CheckSession(string) error
	Delete(key string)
}
type app struct {
	reader     readerService
	parser     parserService
	service    sessionService
	validNames map[string][]byte
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

	return &app{validNames: cfg, parser: parser.New(), reader: reader.New(), service: service.New()}
}
func (a *app) Run() {
	for {
		log.Printf("\nCommands available:\n\t%s %s\n\t%s %s\n\t%s %s\n\t%s\n", COMMAND_LOGIN, ARGS_LOGIN, COMMAND_AUTH, ARGS_AUTH, COMMAND_DELETE, ARGS_DELETE, COMMAND_QUIT)
		cmd, args, err := a.parser.ParseCommand(a.reader.WaitForInput(os.Stdin))
		if err != nil {
			log.Printf("\n! Parse error: %v\n\n", err)
			continue
		}

		switch cmd {
		case COMMAND_QUIT:
			a.Close()
		case COMMAND_LOGIN:
			if len(args) != 2 {
				log.Printf("\n\n! Excepted 2 arguments, but got %d\n\n", len(args))
				continue
			}
			password, ok := a.validNames[args[0]]
			if !ok {
				log.Println("\n\n! User does not exist\n ")
				continue
			}

			if err := bcrypt.CompareHashAndPassword(password, []byte(args[1])); err != nil {
				log.Println("\n\n! Incorrect password\n ")
				continue
			}
			log.Printf("\nGenerated session key: %s\nSession will be expired in 1 minute\n", a.service.CreateSession())
			continue
		case COMMAND_AUTH:
			if len(args) != 1 {
				log.Printf("\n\n! Excepted 1 arguments, but got %d\n\n", len(args))
				continue
			}
			if err := a.service.CheckSession(args[0]); err == nil {
				log.Println("\n\nUser successfully authorized")
				continue
			} else {
				log.Printf("\n\n%v\n", err)
				continue
			}
		case COMMAND_DELETE:
			if len(args) != 1 {
				log.Printf("\n\n! Excepted 1 arguments, but got %d\n\n", len(args))
				continue
			}
			a.service.Delete(args[0])
			log.Println("\n\nSession deleted successfully")
			continue
		default:
			log.Printf("\n\n! Command %s does not exist\n\n", cmd)
			continue
		}
	}
}
func (a *app) Close() {
	os.Exit(0)
}
