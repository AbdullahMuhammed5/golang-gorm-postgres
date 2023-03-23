package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/abdullahmuhammed5/golang-gorm-postgres/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TicketController struct {
	DB *gorm.DB
}

func NewTicketController(DB *gorm.DB) TicketController {
	return TicketController{DB}
}

func (tc *TicketController) CreateTicket(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateTicketRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newTicket := &models.Ticket{
		Title:       payload.Title,
		Description: payload.Description,
		CreatedBy:   currentUser.ID,
		Status:      "NEW",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := tc.DB.Create(&newTicket)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "No ticket found with that ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ticketResponse := &models.TicketResponse{
		ID:          newTicket.ID,
		Title:       newTicket.Title,
		Description: newTicket.Description,
		Owner:       newTicket.Owner,
		Status:      newTicket.Status,
		CreatedAt:   newTicket.CreatedAt,
		UpdatedAt:   newTicket.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": ticketResponse})
}

func (tc *TicketController) UpdateTicket(ctx *gin.Context) {
	ticketId := ctx.Param("ticketId")

	var payload *models.UpdateTicket
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedTicket models.Ticket
	result := tc.DB.First(&updatedTicket, "id = ?", ticketId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No ticket found with that ID"})
		return
	}
	now := time.Now()
	ticketToUpdate := models.Ticket{
		Title:       payload.Title,
		Description: payload.Description,
		CreatedAt:   updatedTicket.CreatedAt,
		UpdatedAt:   now,
	}

	tc.DB.Model(&updatedTicket).Updates(ticketToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedTicket})
}

func (tc *TicketController) FindTicketById(ctx *gin.Context) {
	ticketId := ctx.Param("ticketId")

	var ticket models.Ticket
	result := tc.DB.First(&ticket, "id = ?", ticketId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No ticket found with that ID"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": ticket})
}

func (tc *TicketController) FindTickets(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var tickets []models.Ticket
	results := tc.DB.Limit(intLimit).Offset(offset).Find(&tickets)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(tickets), "data": tickets})
}

func (tc *TicketController) DeleteTicket(ctx *gin.Context) {
	ticketId := ctx.Param("ticketId")

	result := tc.DB.Delete(&models.Ticket{}, "id = ?", ticketId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No ticket found with that ID"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Deleted Successfully"})
}
