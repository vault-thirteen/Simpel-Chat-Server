package rpc

import "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"

type BanUserParams struct {
	Auth   *Auth           `json:"auth,omitempty"`
	UserId common.ObjectId `json:"userId" gorm:"uniqueIndex"`
}

type BanUserResult struct {
	Success
}
