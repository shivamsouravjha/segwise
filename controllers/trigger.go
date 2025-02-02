package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"segwise/clients/postgres"
	"segwise/helpers"
	"segwise/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateTrigger creates a new event trigger
// @Summary Create a new trigger
// @Description Creates a scheduled or API trigger
// @Tags Triggers
// @Accept json
// @Produce json
// @Param trigger body models.Trigger true "Trigger object"
// @Success 201 {object} models.Trigger
// @Failure 400 {object} map[string]string
// @Router /triggers [post]
// @Security BearerAuth
func CreateTrigger(c *gin.Context) {
	db := postgres.GetDB()

	var trigger models.Trigger
	if err := c.ShouldBindJSON(&trigger); err != nil {
		zap.L().Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if trigger.Type != "scheduled" && trigger.Type != "api" {
		zap.L().Error("Trigger type must be 'scheduled' or 'api'")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Trigger type must be 'scheduled' or 'api'"})
		return
	}
	_, err := helpers.ParseOneTimeSchedule(trigger.Schedule)
	if err != nil {
		zap.L().Error("Invalid schedule", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule, please use a schema like 'in 5 minutes'"})
		return
	}
	if err := db.Create(&trigger).Error; err != nil {
		zap.L().Error("Failed to create trigger", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trigger"})
		return
	}

	c.JSON(http.StatusCreated, trigger)

}

// GetTriggers retrieves all triggers
// @Summary Get all triggers
// @Description Fetch all stored triggers
// @Tags Triggers
// @Produce json
// @Success 200 {array} models.Trigger
// @Router /triggers [get]
// @Security BearerAuth
func GetTriggers(c *gin.Context) {
	var triggers []models.Trigger
	db := postgres.GetDB()
	db.Find(&triggers)
	c.JSON(http.StatusOK, triggers)

}

// @Summary Get a specific trigger
// @Description Fetch a trigger by ID
// @Produce  json
// @Param id path string true "Trigger ID"
// @Success 200 {object} models.Trigger
// @Router /triggers/{id} [get]
// @Security BearerAuth
func GetTriggerByID(c *gin.Context) {
	var trigger models.Trigger
	id := c.Param("id")
	db := postgres.GetDB()
	if err := db.First(&trigger, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trigger not found"})
		return
	}

	c.JSON(http.StatusOK, trigger)

}

// @Summary Update a trigger
// @Description Modify an existing trigger
// @Accept  json
// @Produce  json
// @Param id path string true "Trigger ID"
// @Param trigger body models.Trigger true "Updated Trigger object"
// @Success 200 {object} models.Trigger
// @Router /triggers/{id} [put]
// @Security BearerAuth
func UpdateTrigger(c *gin.Context) {

	id := c.Param("id")
	var trigger models.Trigger
	db := postgres.GetDB()

	if err := db.First(&trigger, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trigger not found"})
		return
	}

	if err := c.ShouldBindJSON(&trigger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&trigger)
	c.JSON(http.StatusOK, trigger)

}

// @Summary Delete a trigger
// @Description Removes a trigger from the system
// @Param id path string true "Trigger ID"
// @Success 200 {string} string "Trigger deleted"
// @Router /triggers/{id} [delete]
// @Security BearerAuth
func DeleteTrigger(c *gin.Context) {
	db := postgres.GetDB()

	id := c.Param("id")
	if err := db.Delete(&models.Trigger{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trigger"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Trigger deleted"})

}

type TriggerRequest struct {
	Endpoint string            `json:"endpoint"`
	Payload  map[string]string `json:"payload"`
}

// TestAPITrigger
// @Summary Test a one-time API trigger
// @Description This endpoint sends an API request with a **test payload** to a specified endpoint without saving it as a trigger.
// @Tags Testing API
// @Accept json
// @Produce json
// @Param request body TriggerRequest{Endpoint string `json:"endpoint"`; Payload map[string]string `json:"payload"`} true "API endpoint and payload to test"
// @Success 200 {object} map[string]interface{} "Trigger executed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 500 {object} map[string]interface{} "Error in executing API trigger"
// @Router /triggers/test/api [post]
// @Security BearerAuth
func TestAPITrigger(c *gin.Context) {
	var request TriggerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Simulate API call (without persisting the event)
	client := &http.Client{Timeout: 10 * time.Second}
	jsonPayload, _ := json.Marshal(request.Payload)

	req, err := http.NewRequest("POST", request.Endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute API trigger"})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Test API trigger executed", "response_status": resp.StatusCode})

}
