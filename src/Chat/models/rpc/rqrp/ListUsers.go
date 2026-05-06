package rqrp

import (
	usr "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/User"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/rpc"
)

type ListUsersParams struct {
	Auth       *rpc.Auth `json:"auth,omitempty"`
	PageSize   int       `json:"pageSize"`
	PageNumber int       `json:"pageNumber"`
}

type ListUsersResult struct {
	PageSize   int          `json:"pageSize"`
	PageNumber int          `json:"pageNumber"`
	TotalPages int          `json:"totalPages"`
	TotalItems int          `json:"totalItems"`
	Items      int          `json:"items"`
	Users      []*usr.User2 `json:"users"`
}
