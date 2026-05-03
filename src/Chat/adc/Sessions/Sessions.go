package ss

import (
	"errors"
	"sync"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/database"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/der"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/generator"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	ev "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/Event"
	ses "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/Session"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
	re "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/rpc/errors"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type Sessions struct {
	guard              sync.RWMutex
	db                 *database.Database
	der                *der.DatabaseErrorReporter
	criticalErrorsChan *chan error
	generator          *generator.Generator

	sessionCountMax          int
	inactivityDurationMaxSec int64
	count                    int
	sessionsByUserId         map[common.ObjectId]*ses.Session
	sessionsByToken          map[string]*ses.Session
}

func NewSessions(
	db *database.Database,
	der *der.DatabaseErrorReporter,
	criticalErrorsChan *chan error,
	generator *generator.Generator,
	sessionCountMax int,
	inactivityDurationMaxSec int64,
) (s *Sessions, err error) {
	s = &Sessions{
		guard:              sync.RWMutex{},
		db:                 db,
		der:                der,
		criticalErrorsChan: criticalErrorsChan,
		generator:          generator,

		sessionCountMax:          sessionCountMax,
		inactivityDurationMaxSec: inactivityDurationMaxSec,
		count:                    0,
		sessionsByUserId:         make(map[common.ObjectId]*ses.Session),
		sessionsByToken:          make(map[string]*ses.Session),
	}

	err = s.resetSessionsInDb()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Sessions) StartSession(userId common.ObjectId) (token *string, rpcErr *jrm1.RpcError) {
	s.guard.Lock()
	defer s.guard.Unlock()

	var err error
	token, err = s.generator.TG().CreatePassword()
	if err != nil {
		return nil, re.NewRpcError_TokenGenerator(err)
	}

	var session = ses.NewSession(userId, token)

	// Check for duplicates.
	{
		_, sessionAlreadyExists := s.sessionsByUserId[session.UserId]
		if sessionAlreadyExists {
			return nil, re.NewRpcError_SessionAlreadyExists(nil)
		}

		_, sessionAlreadyExists = s.sessionsByToken[*token]
		if sessionAlreadyExists {
			return nil, re.NewRpcError_SessionAlreadyExists(nil)
		}
	}

	// Check limits.
	if s.count >= s.sessionCountMax {
		return nil, re.NewRpcError_SessionCountLimit(nil)
	}

	rpcErr = s.addSession(session, *token)
	if rpcErr != nil {
		return nil, rpcErr
	}

	rpcErr = s.internalSelfCheck()
	if rpcErr != nil {
		return nil, rpcErr
	}

	return token, nil
}

// GetUserSession function tries to get an active session of a user.
// If a session is outdated, e.g. when a user is inactive for too long, it
// automatically logs the user out.
func (s *Sessions) GetUserSession(token string) (session *ses.Session, rpcErr *jrm1.RpcError) {
	s.guard.Lock()
	defer s.guard.Unlock()

	var sessionExists bool
	session, sessionExists = s.sessionsByToken[token]
	if !sessionExists {
		return nil, re.NewRpcError_SessionIsNotFound(nil)
	}

	inactiveTime := time.Now().UTC().Unix() - session.GetLastActivityTimeTS()
	if inactiveTime > s.inactivityDurationMaxSec {
		// Session has timed out. Log the user out.
		rpcErr = s.stopExistingSession(session, token)
		if rpcErr != nil {
			return nil, rpcErr
		}

		// Event report.
		{
			event := &ev.Event{
				ActorId:      session.UserId,
				Type:         enum.EventType_SessionTimeOut,
				TargetUserId: nil,
			}

			err := s.db.CreateEvent(event)
			if err != nil {
				return nil, s.der.DatabaseError(err)
			}
		}

		return nil, re.NewRpcError_SessionHasTimedOut(nil)
	}

	return session, nil
}
func (s *Sessions) StopExistingSession(session *ses.Session, token string) (rpcErr *jrm1.RpcError) {
	s.guard.Lock()
	defer s.guard.Unlock()
	return s.stopExistingSession(session, token)
}
func (s *Sessions) StopUserSessionIfExists(userId common.ObjectId) (sessionWasFound bool, rpcErr *jrm1.RpcError) {
	s.guard.Lock()
	defer s.guard.Unlock()

	session, sessionExists := s.sessionsByUserId[userId]
	if !sessionExists {
		return false, nil
	}

	if (session.UserId != userId) || (session.Token == nil) {
		return true, re.NewRpcError_ActiveDataController(nil)
	}

	rpcErr = s.stopExistingSession(session, *session.Token)
	if rpcErr != nil {
		return true, rpcErr
	}

	return true, nil
}
func (s *Sessions) RemoveOutdatedSessions() (rpcErr *jrm1.RpcError) {
	s.guard.Lock()
	defer s.guard.Unlock()

	var inactiveTime int64
	for token, session := range s.sessionsByToken {
		inactiveTime = time.Now().UTC().Unix() - session.GetLastActivityTimeTS()
		if inactiveTime > s.inactivityDurationMaxSec {
			rpcErr = s.stopExistingSession(session, token)
			if rpcErr != nil {
				return rpcErr
			}

			// Event report.
			{
				event := &ev.Event{
					ActorId:      session.UserId,
					Type:         enum.EventType_SessionTimeOut,
					TargetUserId: nil,
				}

				err := s.db.CreateEvent(event)
				if err != nil {
					return s.der.DatabaseError(err)
				}
			}
		}
	}

	return nil
}

func (s *Sessions) resetSessionsInDb() (err error) {
	err = s.db.CleanAllSessions()
	if err != nil {
		return err
	}

	return nil
}
func (s *Sessions) addSession(session *ses.Session, token string) (rpcErr *jrm1.RpcError) {
	err := s.db.CreateSession(session)
	if err != nil {
		return s.der.DatabaseError(err)
	}

	s.sessionsByUserId[session.UserId] = session
	s.sessionsByToken[token] = session
	s.count++

	return nil
}
func (s *Sessions) deleteSession(session *ses.Session, token string) (rpcErr *jrm1.RpcError) {
	err := s.db.DeleteSession(session)
	if err != nil {
		return s.der.DatabaseError(err)
	}

	delete(s.sessionsByUserId, session.UserId)
	delete(s.sessionsByToken, token)
	s.count--

	return nil
}
func (s *Sessions) internalSelfCheck() (rpcErr *jrm1.RpcError) {
	// 1. Compare session count in memory.
	if (s.count != len(s.sessionsByUserId)) ||
		(s.count != len(s.sessionsByToken)) {
		defer func() {
			*(s.criticalErrorsChan) <- helper.NewError_WrappedError(helper.Err_ADC, errors.New(helper.Err_SessionsCountMismatch))
		}()

		return re.NewRpcError_ActiveDataController(nil)
	}

	// 2. Compare session count in database and memory.
	sessionsCountInDb, err := s.db.CountAllSessions()
	if err != nil {
		return s.der.DatabaseError(err)
	}

	if sessionsCountInDb != s.count {
		defer func() {
			*(s.criticalErrorsChan) <- helper.NewError_WrappedError(helper.Err_ADC, errors.New(helper.Err_SessionsCountMismatch))
		}()

		return re.NewRpcError_ActiveDataController(nil)
	}

	return nil
}
func (s *Sessions) stopExistingSession(session *ses.Session, token string) (rpcErr *jrm1.RpcError) {
	rpcErr = s.deleteSession(session, token)
	if rpcErr != nil {
		return rpcErr
	}

	rpcErr = s.internalSelfCheck()
	if rpcErr != nil {
		return rpcErr
	}

	return nil
}
