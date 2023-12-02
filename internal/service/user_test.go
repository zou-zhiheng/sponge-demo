package service

import (
	"fmt"
	"testing"
	"user/configs"
	"user/internal/config"
	"user/internal/model"
)

func init() {
	_ = config.Init(configs.Path("user.docker-compose-yaml"))
	model.InitGoMybatis()
}

func TestGetUserById(t *testing.T) {
	userMapper := model.GetUserMapper()
	res, err := userMapper.SelectById(1714934048536662016)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}

func TestGenerateSymmetricKey(t *testing.T) {

}
