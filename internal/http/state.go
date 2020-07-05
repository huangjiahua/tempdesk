package http

import (
	td "github.com/huangjiahua/tempdesk"
	"github.com/huangjiahua/tempdesk/internal/auth"
	"net/http"
)

type State struct {
	Users  td.UserService
	Files  td.FileService
	Auther auth.UserAuther
}

func (s *State) AuthUser(req *http.Request) (td.User, error) {
	return s.Auther.AuthUser(req, s.Users)
}
