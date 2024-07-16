package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poligonoio/vega-core/internal/services"
)

type ProbeController struct {
	ProbeService services.ProbeService
}

func NewProbeController(probeService services.ProbeService) ProbeController {
	return ProbeController{
		ProbeService: probeService,
	}
}

func (self *ProbeController) ReadinessProbe(c *gin.Context) {
	// placeholder
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (self *ProbeController) LivenessProbe(c *gin.Context) {
	// placeholder
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (self *ProbeController) RegisterProbeRoutes(rg *gin.RouterGroup) {
	rg.GET("/readiness", self.ReadinessProbe)
	rg.GET("/liveness", self.LivenessProbe)
}
