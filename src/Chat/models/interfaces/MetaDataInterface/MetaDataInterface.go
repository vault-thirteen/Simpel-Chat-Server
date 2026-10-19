package mdi

import "github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"

type MetaDataInterface interface {
	GetMetaData() (md *common.MetaData)
}
