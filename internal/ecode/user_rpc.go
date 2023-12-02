package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// user business-level rpc error codes.
// the _userNO value range is 1~100, if the same number appears, it will cause a failure to start the service.
var (
	_userNO       = 15
	_userName     = "user"
	_userBaseCode = errcode.RCode(_userNO)

	StatusCreateUser         = errcode.NewRPCStatus(_userBaseCode+1, "failed to create "+_userName)
	StatusDeleteByIDUser     = errcode.NewRPCStatus(_userBaseCode+2, "failed to delete "+_userName)
	StatusDeleteByIDsUser    = errcode.NewRPCStatus(_userBaseCode+3, "failed to delete by batch ids "+_userName)
	StatusUpdateByIDUser     = errcode.NewRPCStatus(_userBaseCode+4, "failed to update "+_userName)
	StatusGetByIDUser        = errcode.NewRPCStatus(_userBaseCode+5, "failed to get "+_userName+" details")
	StatusGetByConditionUser = errcode.NewRPCStatus(_userBaseCode+6, "failed to get "+_userName+" by conditions")
	StatusListByIDsUser      = errcode.NewRPCStatus(_userBaseCode+7, "failed to list by batch ids "+_userName)
	StatusListUser           = errcode.NewRPCStatus(_userBaseCode+8, "failed to list of "+_userName)
	// error codes are globally unique, adding 1 to the previous error code
)
