package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/reversersed/zyfra/internal/handlers"
	"github.com/reversersed/zyfra/internal/handlers/models"
	"github.com/reversersed/zyfra/internal/service"
	"github.com/reversersed/zyfra/internal/storage"
	"github.com/reversersed/zyfra/pkg/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var engine *gin.Engine

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	var mongoContainer testcontainers.Container
	for i := 0; i < 5; i++ {
		req := testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForListeningPort("27017/tcp"),
			Env: map[string]string{
				"MONGO_INITDB_ROOT_USERNAME": "root",
				"MONGO_INITDB_ROOT_PASSWORD": "root",
			},
		}
		mongoContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err == nil {
			break
		}
		log.Printf("failed to create container: %v, retry %d/5", err, i+1)
		<-time.After(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not start mongo: %s", err)
	}
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop mongo: %s", err)
		}
	}()
	host, err := mongoContainer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}
	port, err := mongoContainer.MappedPort(ctx, "27017/tcp")
	if err != nil {
		log.Fatal(err)
	}
	cfg := &mongo.DatabaseConfig{
		Host:     host,
		Port:     port.Int(),
		User:     "root",
		Password: "root",
		Base:     "testbase",
		AuthDb:   "admin",
	}

	driver, err := mongo.NewClient(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}
	storage := storage.New(driver)
	service := service.New(storage)
	h := handlers.New(service)
	engine = gin.Default()
	gin.SetMode(gin.TestMode)
	h.Register(engine)

	os.Exit(m.Run())
}

type testCase struct {
	Name             string
	ExceptedCode     int
	ExceptedResponse string
	RouteAddition    string
	Request          any
}
type testTable struct {
	Method string
	Route  string
	Cases  []testCase
}

func proceedTest(t *testing.T, table testTable) {
	for _, c := range table.Cases {
		t.Run(c.Name, func(t *testing.T) {

			body := []byte{}
			if c.Request != nil {
				body, _ = json.Marshal(c.Request)
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(table.Method, table.Route+c.RouteAddition, bytes.NewReader(body))
			engine.ServeHTTP(w, r)

			if len(c.ExceptedResponse) != 0 {
				response, _ := io.ReadAll(w.Body)
				assert.Equal(t, c.ExceptedResponse, string(response))
			}
			assert.Equal(t, c.ExceptedCode, w.Result().StatusCode)
		})
	}
}

func TestLoginHandler(t *testing.T) {
	table := testTable{
		Method: http.MethodPost,
		Route:  "/api/sessions",
		Cases: []testCase{
			{
				Name:             "empty request",
				ExceptedCode:     http.StatusBadRequest,
				Request:          nil,
				ExceptedResponse: "{\"message\":\"Excepted non-empty login\",\"error\":\"login length was zero\"}",
			},
			{
				Name:             "empty password",
				ExceptedCode:     http.StatusBadRequest,
				Request:          models.LoginCommand{Login: "admin"},
				ExceptedResponse: "{\"message\":\"Excepted non-empty password\",\"error\":\"password length was zero\"}",
			},
			{
				Name:         "successful login",
				ExceptedCode: http.StatusOK,
				Request:      models.LoginCommand{Login: "admin", Password: "admin"},
			},
			{
				Name:             "wrong password",
				ExceptedCode:     http.StatusUnauthorized,
				Request:          models.LoginCommand{Login: "admin", Password: "ad"},
				ExceptedResponse: "{\"message\":\"User not found\",\"error\":\"crypto/bcrypt: hashedPassword is not the hash of the given password\"}",
			},
			{
				Name:             "unexisting user",
				ExceptedCode:     http.StatusUnauthorized,
				Request:          models.LoginCommand{Login: "user", Password: "ad"},
				ExceptedResponse: "{\"message\":\"User not found\",\"error\":\"mongo: no documents in result\"}",
			},
		},
	}
	proceedTest(t, table)
}
func TestAuthHandler(t *testing.T) {
	body, _ := json.Marshal(&models.LoginCommand{Login: "admin", Password: "admin"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/sessions", bytes.NewReader(body))
	engine.ServeHTTP(w, r)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
	response, _ := io.ReadAll(w.Body)
	var session struct {
		Session string
	}
	json.Unmarshal(response, &session)

	table := testTable{
		Method: http.MethodGet,
		Route:  "/api/sessions",
		Cases: []testCase{
			{
				Name:          "successful authorization",
				RouteAddition: "/" + session.Session,
				ExceptedCode:  http.StatusOK,
			},
			{
				Name:             "invalid session key",
				RouteAddition:    "/dsa",
				ExceptedCode:     http.StatusUnauthorized,
				ExceptedResponse: "{\"message\":\"User not authorized\",\"error\":\"Session not found\"}",
			},
		},
	}
	proceedTest(t, table)
}

func TestDeleteHandler(t *testing.T) {
	body, _ := json.Marshal(&models.LoginCommand{Login: "admin", Password: "admin"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/sessions", bytes.NewReader(body))
	engine.ServeHTTP(w, r)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)
	response, _ := io.ReadAll(w.Body)
	var session struct {
		Session string
	}
	json.Unmarshal(response, &session)

	table := testTable{
		Method: http.MethodDelete,
		Route:  "/api/sessions",
		Cases: []testCase{
			{
				Name:          "successful deletion",
				RouteAddition: "/" + session.Session,
				ExceptedCode:  http.StatusOK,
			},
			{
				Name:             "second deletion not found",
				RouteAddition:    "/" + session.Session,
				ExceptedCode:     http.StatusNotFound,
				ExceptedResponse: "{\"message\":\"Session was not deleted\",\"error\":\"session not found\"}",
			},
			{
				Name:             "invalid session key",
				RouteAddition:    "/invalidKey",
				ExceptedCode:     http.StatusNotFound,
				ExceptedResponse: "{\"message\":\"Session was not deleted\",\"error\":\"the provided hex string is not a valid ObjectID\"}",
			},
		},
	}
	proceedTest(t, table)
}
