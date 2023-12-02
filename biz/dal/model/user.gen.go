// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID         int64     `gorm:"column:id;primaryKey" json:"id,string,omitempty" form:"id"`
	Name       string    `gorm:"column:name;not null" json:"name" form:"name"`
	Account    string    `gorm:"column:account;not null" json:"account" form:"account"`
	Password   string    `gorm:"column:password;not null" json:"password" form:"password"`
	Salt       string    `gorm:"column:salt;not null" json:"salt" form:"salt"`
	Sex        string    `gorm:"column:sex" json:"sex" form:"sex"`
	Age        int16     `gorm:"column:age" json:"age" form:"age"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time,omitempty" form:"update_time"`
	Valid      int16     `gorm:"column:valid;not null" json:"valid" form:"valid"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
