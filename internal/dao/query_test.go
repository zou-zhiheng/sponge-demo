package dao

import (
	"fmt"
	"testing"
	"user/internal/dao/query"
	"user/internal/model"
)

func TestCreate(t *testing.T) {
	err := query.Use(model.GetDB()).User.Create(&model.User{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Create Pass")
}

func TestSelect(t *testing.T) {
	//u := query.User
	//u.UseDB(model.GetDB())
	//user, err := u.WithContext(context.Background()).Where(u.Name.Eq("")).Find()
}
