package rqrp

import (
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type ChangePassword2Params struct {
	Auth             *rpc.Auth `json:"auth"`
	RequestId        string    `json:"requestId"`
	VerificationCode string    `json:"verificationCode"`
	UserPassword     string    `json:"userPassword"`
	NewUserPassword1 string    `json:"newUserPassword1"`
	NewUserPassword2 string    `json:"newUserPassword2"`
}

type ChangePassword2Result struct{}
