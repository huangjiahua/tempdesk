package mock

import (
	td "github.com/huangjiahua/tempdesk"
)

type UserService struct {
	M map[string]td.User
}

func NewUserService() *UserService {
	return &UserService{make(map[string]td.User)}
}

func (u *UserService) CreateUser(user td.User) (err error) {
	if _, ok := u.M[user.Name]; ok {
		err = &td.UserServiceError{Kind: td.NameAlreadyExists, Err: nil}
		return
	}
	u.M[user.Name] = user
	return
}

func (u *UserService) UpdateUser(user td.User) (err error) {
	if _, ok := u.M[user.Name]; !ok {
		err = &td.UserServiceError{Kind: td.NameNotExists, Err: nil}
		return
	}
	u.M[user.Name] = user
	return
}

func (u *UserService) User(name string) (user td.User, ok bool) {
	user, ok = u.M[name]
	return
}
