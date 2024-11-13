package handlers

import "github.com/gin-gonic/gin"

//go:generate mockgen -source=module.go -destination=mocks/module.go

type Handler interface {
	Register(*gin.Engine)
}
type sessionService interface {
	CreateSession() string
	CheckSession(string) error
	Delete(key string) error
}
type handler struct {
	service sessionService
	cfg     map[string][]byte
}

func New(service sessionService, cfg map[string][]byte) Handler {
	return &handler{service: service, cfg: cfg}
}
func (h *handler) Register(g *gin.Engine) {
	group := g.Group("/api")
	{
		group.GET("/sessions/:session", h.HandleAuthRequest)
		group.DELETE("/sessions/:session", h.HandleDeleteCommand)
		group.POST("/sessions", h.HandleLoginCommand)
	}
}
