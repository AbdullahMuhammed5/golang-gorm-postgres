package routes

import (
	"github.com/abdullahmuhammed5/golang-gorm-postgres/controllers"
	"github.com/abdullahmuhammed5/golang-gorm-postgres/middleware"
	"github.com/gin-gonic/gin"
)

type TicketRouteController struct {
	ticketController controllers.TicketController
}

func NewRouteTicketController(ticketController controllers.TicketController) TicketRouteController {
	return TicketRouteController{ticketController}
}

func (pc *TicketRouteController) TicketRoute(rg *gin.RouterGroup) {

	router := rg.Group("tickets")
	router.Use(middleware.DeserializeUser())
	router.POST("/", pc.ticketController.CreateTicket)
	router.GET("/", pc.ticketController.FindTickets)
	router.PATCH("/:ticketId", pc.ticketController.UpdateTicket)
	router.GET("/:ticketId", pc.ticketController.FindTicketById)
	router.DELETE("/:ticketId", pc.ticketController.DeleteTicket)
}
