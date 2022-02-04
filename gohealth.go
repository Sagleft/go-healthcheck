package gohealth

const (
	defaultPort = 8080
)

// HandlerTask - healtch check handler constructor data
type HandlerTask struct {
	DisableGin bool // disable web-server
	ListenPort int
}

// Handler - go-healthcheck handler
type Handler struct {
	task        HandlerTask
	checkpoints []Checkpoint
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
func NewHandler(task HandlerTask) *Handler {
	return &Handler{
		task:        task,
		checkpoints: make([]Checkpoint, 0),
	}
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

// AddCheckpoint - add new health checkpoint
func (h *Handler) AddCheckpoint(data CheckpointData) {
	h.checkpoints = append(h.checkpoints, newCheckpoint(data))
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
