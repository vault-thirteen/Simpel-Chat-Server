package der

import (
	"log"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"

	re "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/rpc/errors"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type DatabaseErrorReporter struct {
	criticalErrorsChan *chan error
}

func NewDatabaseErrorReporter(criticalErrorsChan *chan error) *DatabaseErrorReporter {
	return &DatabaseErrorReporter{
		criticalErrorsChan: criticalErrorsChan,
	}
}

func (der *DatabaseErrorReporter) DatabaseError(err error) *jrm1.RpcError {
	der.processDatabaseError(err)
	return re.NewRpcError_DatabaseError(err)
}

func (der *DatabaseErrorReporter) processDatabaseError(err error) {
	if err == nil {
		return
	}

	if helper.IsNetworkError(err) {
		*(der.criticalErrorsChan) <- err
		return
	}

	der.logError(err)
	return
}

func (der *DatabaseErrorReporter) logError(err error) {
	if err == nil {
		return
	}

	log.Println(err)
}
