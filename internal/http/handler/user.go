package handler

import (
	"bytes"
	"encoding/json"
	thttp "github.com/huangjiahua/tempdesk/internal/http"
	"io"
	"net/http"
)

type User struct {
	state *thttp.State
}

func (u *User) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(res, "Error parsing url arguments or form", http.StatusBadRequest)
		return
	}
	user, err := u.state.AuthUser(req)
	if err != nil {
		http.Error(res, "Error authenticating", http.StatusForbidden)
		return
	}

	ret := map[string]string{"name": user.Name}
	body, err := json.Marshal(ret)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}
	r := bytes.NewReader(body)
	_, err = io.Copy(res, r)
	if err != nil {
		// log
	}
}
