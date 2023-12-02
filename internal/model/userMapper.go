package model

type UserMapper struct {
	SelectAll  func() ([]User, error)
	SelectById func(id uint64) (User, error)  `args:"id"`
	UpdateById func(user User) (int64, error) `args:"args"`
	Insert     func(user User) (int64, error) `args:"args"`
	DeleteById func(id int64) (int64, error)  `args:"id"`
	JoinTest   func() ([]map[string]interface{}, error)
}
