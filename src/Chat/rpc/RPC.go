package rpc

import (
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/adc"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/database"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/der"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/generator"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/mailer"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/settings"
)

const (
	RpcDurationFieldName  = "dur"
	RpcRequestIdFieldName = "rid"
)

type RPC struct {
	processor  *jrm1.Processor
	controller *RpcController
}

func NewRPC(
	db *database.Database,
	mailer *mailer.Mailer,
	generator *generator.Generator,
	adc *adc.ActiveDataController,
	der *der.DatabaseErrorReporter,
	chatUserSettings *settings.ChatUserSettings,
) (rpc *RPC, err error) {
	rpc = new(RPC)

	// Controller.
	rpc.controller = NewRpcController(db, mailer, generator, adc, der, chatUserSettings)

	// Processor.
	{
		fnDur := RpcDurationFieldName
		fnReqId := RpcRequestIdFieldName

		ps := &jrm1.ProcessorSettings{
			CatchExceptions:    true,
			LogExceptions:      true,
			CountRequests:      true,
			DurationFieldName:  &fnDur,
			RequestIdFieldName: &fnReqId,
		}

		rpc.processor, err = jrm1.NewProcessor(ps)
		if err != nil {
			return nil, err
		}

		funcs := rpc.controller.GetRpcFunctions()

		for _, fn := range funcs {
			err = rpc.processor.AddFunc(fn)
			if err != nil {
				return nil, err
			}
		}
	}

	return rpc, nil
}

func (rpc *RPC) GetProcessor() *jrm1.Processor { return rpc.processor }
