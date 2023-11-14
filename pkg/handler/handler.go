package handler

import (
	"L0/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("templates/*.html")
	order := router.Group("/order")
	{
		order.GET("/:uid", h.getByUID)
	}

	return router

}
