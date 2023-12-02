package dao

import "gorm.io/gen"

type Querier interface {
	//select *from @@table
	//{{where}}
	//  {{for _,user:=range user}}
	//    {{if user.Name != "" && user.Age >0 }}
	//      (name = @user.Name and age = @user.Age and account like concat("%",@user.account,"%"))
	//    {{end}}
	//   {{end}}
	//{{end}}
	FilterFor(users []gen.T) ([]gen.T, error)
}
