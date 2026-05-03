package cleaner

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/adc"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/database"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/request"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/settings"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type Cleaner struct {
	settings           *CleanerSettings
	db                 *database.Database
	criticalErrorsChan *chan error
	wg                 *sync.WaitGroup
	mustStop           *atomic.Bool
	funcs              []common.ScheduledFn
	adc                *adc.ActiveDataController
}

func NewCleaner(stn *settings.OtherChatSettings, db *database.Database, criticalErrorsChan *chan error, adc *adc.ActiveDataController) *Cleaner {
	c := &Cleaner{
		settings:           NewCleanerSettings(stn),
		db:                 db,
		criticalErrorsChan: criticalErrorsChan,
		wg:                 new(sync.WaitGroup),
		adc:                adc,
	}

	c.mustStop = new(atomic.Bool)
	c.mustStop.Store(false)

	c.funcs = []common.ScheduledFn{
		c.removeOutdatedRegistrationRequests,
		c.removeOutdatedLogInRequests,
		c.removeOutdatedLogOutRequests,
		c.removeOutdatedSessions,
	}

	return c
}

func (c *Cleaner) removeOutdatedRegistrationRequests() (err error) {
	edgeTime := time.Now().UTC().Add(-time.Duration(c.settings.RegistrationRequestTtl) * time.Second)

	var rrs []rq.Registration
	for {
		rrs, err = c.db.GetFirstOutdatedRegistrationRequest(edgeTime)
		if err != nil {
			return err
		}

		if len(rrs) < 1 {
			break
		}

		err = c.db.DeleteRegistrationRequest(&rrs[0])
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Cleaner) removeOutdatedLogInRequests() (err error) {
	edgeTime := time.Now().UTC().Add(-time.Duration(c.settings.RegistrationRequestTtl) * time.Second)

	var lirs []rq.LogIn
	for {
		lirs, err = c.db.GetFirstOutdatedLogInRequest(edgeTime)
		if err != nil {
			return err
		}

		if len(lirs) < 1 {
			break
		}

		err = c.db.DeleteLogInRequest(&lirs[0])
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Cleaner) removeOutdatedLogOutRequests() (err error) {
	edgeTime := time.Now().UTC().Add(-time.Duration(c.settings.RegistrationRequestTtl) * time.Second)

	var lors []rq.LogOut
	for {
		lors, err = c.db.GetFirstOutdatedLogOutRequest(edgeTime)
		if err != nil {
			return err
		}

		if len(lors) < 1 {
			break
		}

		err = c.db.DeleteLogOutRequest(&lors[0])
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Cleaner) removeOutdatedPasswordChangeRequests() (err error) {
	edgeTime := time.Now().UTC().Add(-time.Duration(c.settings.RegistrationRequestTtl) * time.Second)

	var pcrs []rq.ChangePassword
	for {
		pcrs, err = c.db.GetFirstOutdatedPasswordChangeRequest(edgeTime)
		if err != nil {
			return err
		}

		if len(pcrs) < 1 {
			break
		}

		err = c.db.DeletePasswordChangeRequest(&pcrs[0])
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *Cleaner) removeOutdatedSessions() (err error) {
	var rpcErr *jrm1.RpcError
	rpcErr = c.adc.RemoveOutdatedSessions()
	if rpcErr != nil {
		return rpcErr.AsError()
	}

	return nil
}

func (c *Cleaner) Start() {
	c.wg.Add(1)
	go c.run()
}

// Async.
func (c *Cleaner) run() {
	log.Println(helper.Msg_ChatCleanerHasStarted)

	defer func() {
		log.Println(helper.Msg_ChatCleanerHasStopped)
		c.wg.Done()
	}()

	// Time counter.
	// It counts seconds and resets every 24 hours.
	var tc uint = 1
	const SecondsInDay = 86400 // 60*60*24.
	var err error

	for {
		if c.mustStop.Load() {
			break
		}

		// Periodical tasks (every minute).
		if tc%60 == 0 {
			for _, fn := range c.funcs {
				err = fn()
				if err != nil {
					*c.criticalErrorsChan <- err
					return
				}
			}
		}

		// Next tick.
		if tc == SecondsInDay {
			tc = 0
		}
		tc++
		time.Sleep(time.Second)
	}
}

func (c *Cleaner) AskToStop() {
	c.mustStop.Store(true)
}

func (c *Cleaner) WaitForStop() {
	c.wg.Wait()
}
