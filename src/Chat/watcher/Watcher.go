package watcher

import (
	"log"
	"sync"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/controls"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type Watcher struct {
	controls *controls.Controls
	wg       *sync.WaitGroup
	stopChan *chan bool
}

func NewWatcher(controls *controls.Controls) *Watcher {
	w := &Watcher{
		controls: controls,
		wg:       new(sync.WaitGroup),
	}

	// We need a buffered channel, because we need to send a stop signal to
	// the watcher in the 'Stop' method. There is non-zero possibility that
	// the watcher has already stopped reading its channel when the 'Stop'
	// method is running during the emergency shutdown procedure.
	sc := make(chan bool, 1)
	w.stopChan = &sc

	return w
}

func (w *Watcher) Start() {
	w.wg.Add(1)
	go w.run()
}

// Async.
func (w *Watcher) run() {
	log.Println(helper.Msg_ChatWatcherHasStarted)

	// In order not to spam uncontrollable goroutines, this function needs a
	// WaitGroup guard for it to return in case when nothing bad happens.
	defer func() {
		log.Println(helper.Msg_ChatWatcherHasStopped)
		w.wg.Done()
	}()

	var criticalError error

	select {
	case <-*w.stopChan:
		return
	case criticalError = <-*w.controls.GetCriticalErrorsChan():
		log.Println(helper.NewError_GenericError(helper.Err_Critical, criticalError.Error()))
		go w.controls.GetEmergencyShutdownFn()()
	}
}

func (w *Watcher) AskToStop() {
	*w.stopChan <- true
}

func (w *Watcher) WaitForStop() {
	w.wg.Wait()
}
