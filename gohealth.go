package gohealth

// Handler - go-healthcheck handler
type Handler struct {
	Checkpoints []Checkpoint
}

// Signal - service health signal
type Signal struct {
	CheckPassed bool
	ErrorInfo   string
}

type healthCheckCallback func() Signal

// Checkpoint : service must pass all checkpoints before it comes to the boss battle!
type Checkpoint struct {
	name     string
	callback healthCheckCallback
}

// NewHandler - create new health check handler for service
func NewHandler() *Handler {
	return &Handler{}
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
	CheckCallback healthCheckCallback
}

// HealthCheckAnalytics - the result of the service health analysis
type HealthCheckAnalytics struct{}

// AddCheckpoint - add new health checkpoint
func (h *Handler) AddCheckpoint(data CheckpointData) {
	h.Checkpoints = append(h.Checkpoints, newCheckpoint(data))
}

// Check service
func (h *Handler) Check() []Signal {
	indicators := []Signal{}
	for _, checkpoint := range h.Checkpoints {
		indicators = append(indicators, checkpoint.callback())
	}
	return indicators
}
