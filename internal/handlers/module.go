package handlers

import "github.com/gin-gonic/gin"

type Handler interface {
	Register(*gin.Engine)
}
type sessionService interface {
	CreateSession(string, string) (string, error)
	CheckSession(string) error
	Delete(key string) error
}
type handler struct {
	service sessionService
}

func New(service sessionService) Handler {
	return &handler{service: service}
}
func (h *handler) Register(g *gin.Engine) {
	group := g.Group("/api")
	{
		group.GET("/sessions/:session", h.HandleAuthRequest)
		group.DELETE("/sessions/:session", h.HandleDeleteCommand)
		group.POST("/sessions", h.HandleLoginCommand)
	}
}
