package auth

import (
	td "github.com/huangjiahua/tempdesk"
	"net/http"
)

const (
	NotAuthed      string = "Not Authed"
	NoUser         string = "No Such User"
	WrongFormat    string = "Wrong Format"
	Outdated       string = "Message Outdated"
	AutherInternal string = "Auther Internal Error"
)

type UserAuther interface {
	AuthUser(req *http.Request, us td.UserService) (td.User, error)
}

type AutherError struct {
	Kind   string
	Detail string
}

func (a *AutherError) Error() string {
	return a.Kind + ": " + a.Detail
}
