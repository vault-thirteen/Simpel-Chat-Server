package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type PingParams = struct{}

type PingResult struct {
	rpc.Success
}
