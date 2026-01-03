package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gonotifysys.com/internal/worker"
)

type NotificationHandler struct {
	WP *worker.WorkerPool
}

func NewNotificationHandler(wp *worker.WorkerPool) *NotificationHandler {
	return &NotificationHandler{
		WP: wp,
	}
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {

	var req struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	jobID := fmt.Sprintf("%d", time.Now().UnixNano())

	job := worker.Job{
		ID:   jobID,
		Type: req.Type,
	}

	select {
	case h.WP.JobQueue <- job:
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Request queued",
			"job_id":  jobID,
		})
	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Server is busy, too many request",
		})
	}

}
