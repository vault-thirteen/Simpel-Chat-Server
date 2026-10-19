package controls

import (
	"sync/atomic"
)

type Controls struct {
	criticalErrorsChan  *chan error
	chatStoppedChan     *chan bool
	isEmergencyShutdown *atomic.Bool
	emergencyShutdownFn func()
}

func NewControls(emergencyShutdownFn func()) (c *Controls) {
	c = new(Controls)

	cec := make(chan error, 8)
	c.criticalErrorsChan = &cec

	// Channel for external viewers showing that chat has unexpectedly
	// stopped working.
	csc := make(chan bool, 1)
	c.chatStoppedChan = &csc

	c.isEmergencyShutdown = new(atomic.Bool)
	c.isEmergencyShutdown.Store(false)

	c.emergencyShutdownFn = emergencyShutdownFn

	return c
}

func (c *Controls) GetCriticalErrorsChan() *chan error {
	return c.criticalErrorsChan
}
func (c *Controls) GetChatStoppedChan() *chan bool {
	return c.chatStoppedChan
}
func (c *Controls) GetIsEmergencyShutdown() *atomic.Bool {
	return c.isEmergencyShutdown
}
func (c *Controls) GetEmergencyShutdownFn() func() {
	return c.emergencyShutdownFn
}
