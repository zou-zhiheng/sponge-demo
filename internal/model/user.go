package model

import (
	"time"
)

type User struct {
	ID         uint64    `gorm:"column:id;type:bigint(20);primary_key" json:"id"`
	Name       string    `gorm:"column:name;type:varchar(255);NOT NULL" json:"name"`
	Account    string    `gorm:"column:account;type:varchar(255);NOT NULL" json:"account"`
	Password   string    `gorm:"column:password;type:varchar(255);NOT NULL" json:"password"`
	Salt       string    `gorm:"column:salt;type:varchar(255);NOT NULL" json:"salt"`
	Sex        string    `gorm:"column:sex;type:varchar(10)" json:"sex"`
	Age        int       `gorm:"column:age;type:tinyint(4)" json:"age"`
	CreateTime time.Time `gorm:"column:create_time;type:datetime;NOT NULL" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;type:datetime" json:"updateTime"`
	Valid      int       `gorm:"column:valid;type:tinyint(4);NOT NULL" json:"valid"`
}

// TableName table name
func (m *User) TableName() string {
	return "user"
}
