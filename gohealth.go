package gohealth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort     = "8080"
	defaultResponse = "OK"
)

// HandlerTask - healtch check handler constructor data
type HandlerTask struct {
	DisableGin bool // disable web-server
	DisableLog bool
	ListenPort string
}

// Handler - go-healthcheck handler
type Handler struct {
	task        HandlerTask
	checkpoints []Checkpoint
	gin         *gin.Engine
}

// Signal - service health signal
type Signal struct {
	CheckPassed bool
	ErrorInfo   string
}

// HealthCheckCallback - health check callback
type HealthCheckCallback func() Signal

// Checkpoint : service must pass all checkpoints before it comes to the boss battle!
type Checkpoint struct {
	name     string
	callback HealthCheckCallback
}

// NewHandler - create new health check handler for service
func NewHandler(task HandlerTask) *Handler {
	h := Handler{
		task:        task,
		checkpoints: make([]Checkpoint, 0),
	}
	if task.ListenPort == "" {
		h.task.ListenPort = defaultPort
	}
	if !task.DisableGin {
		gin.SetMode(gin.ReleaseMode)
		h.gin = gin.Default()
		h.setup()
	}
	return &h
}

func (h *Handler) setup() {
	h.setupLog()
	h.gin.GET("/healthcheck", h.doHealthCheck)
	go h.gin.Run(":" + h.task.ListenPort)
}

func (h *Handler) setupLog() {
	h.gin.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		if h.task.DisableLog {
			return ""
		}
		return fmt.Sprintf("[%s] \"%s %s %s %d %s \" %s\"\n",
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	}))
}

func (h *Handler) doHealthCheck(c *gin.Context) {
	for _, checkpoint := range h.checkpoints {
		signalData := checkpoint.callback()
		if !signalData.CheckPassed {
			c.String(http.StatusInternalServerError, signalData.ErrorInfo)
			return
		}
	}
	c.String(http.StatusOK, defaultResponse)
}

func newCheckpoint(data CheckpointData) Checkpoint {
	c := Checkpoint{
		name:     data.Name,
		callback: data.CheckCallback,
	}
	if c.callback == nil {
		c.callback = defaultCallback
	}
	return c
}

func defaultCallback() Signal {
	return Signal{
		CheckPassed: false,
		ErrorInfo:   "callback is not set",
	}
}

// CheckpointData - health check checkpoint data
type CheckpointData struct {
	Name          string
	CheckCallback HealthCheckCallback
}

// AddCheckpoint - add new health checkpoint
func (h *Handler) AddCheckpoint(data CheckpointData) {
	h.checkpoints = append(h.checkpoints, newCheckpoint(data))
}

// AddCheckpoints - add new health checkpoint
func (h *Handler) AddCheckpoints(checkpoints []CheckpointData) {
	for _, data := range checkpoints {
		h.checkpoints = append(h.checkpoints, newCheckpoint(data))
	}
}

// Check service
func (h *Handler) Check() []Signal {
	indicators := []Signal{}
	for _, checkpoint := range h.checkpoints {
		indicators = append(indicators, checkpoint.callback())
	}
	return indicators
}

// SignalError - get new errorsignal
func SignalError(errInfo string) Signal {
	return Signal{
		CheckPassed: false,
		ErrorInfo:   errInfo,
	}
}

// SignalNormal - checkpoint passed
func SignalNormal() Signal {
	return Signal{
		CheckPassed: true,
	}
}
