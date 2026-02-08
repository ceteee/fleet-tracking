package vehicle

import (
	domain "fleet-management-system/internal/fleet/vehicle"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *domain.Service
}

func NewHandler(service *domain.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateLocation(c *gin.Context) {
	var req domain.LocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err := h.service.RecordLocation(
		c.Request.Context(),
		req.VehicleID,
		req.Latitude,
		req.Longitude,
		time.Unix(req.Timestamp, 0),
	)

	if err != nil {
		log.Println("ERROR RecordLocation:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to record location",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "location recorded",
	})
}

func (h *Handler) GetLatestLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "vehicle_id is required",
		})
		return
	}

	loc, err := h.service.GetLatestLocation(
		c.Request.Context(),
		vehicleID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "location not found",
		})
	}

	response := domain.LocationResponse{
		VehicleID: loc.VehicleID,
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
		Timestamp: loc.RecordedAt.Unix(),
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetLocationHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")

	if vehicleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "vehicle_id is required",
		})
		return
	}

	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start and end value are required",
		})
		return
	}

	startUnix, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start timestamp",
		})
		return
	}

	endUnix, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end timestamp",
		})
		return
	}

	locations, err := h.service.GetLocationHistory(
		c.Request.Context(),
		vehicleID,
		time.Unix(startUnix, 0),
		time.Unix(endUnix, 0),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get location history",
		})
		return
	}

	responses := make([]domain.LocationResponse, 0, len(locations))
	for _, loc := range locations {
		responses = append(responses, domain.LocationResponse{
			VehicleID: loc.VehicleID,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
			Timestamp: loc.RecordedAt.Unix(),
		})
	}

	c.JSON(http.StatusOK, responses)
}
