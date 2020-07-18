package handler

import (
	"bytes"
	"encoding/json"
	thttp "github.com/huangjiahua/tempdesk/internal/http"
	tlog "github.com/huangjiahua/tempdesk/internal/log"
	"io"
	"net/http"
)

type User struct {
	state *thttp.State
}

func (u *User) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		tlog.Debug("error parsing url", tlog.String("err", err.Error()))
		http.Error(res, "Error parsing url arguments or form", http.StatusBadRequest)
		return
	}
	user, err := u.state.AuthUser(req)
	if err != nil {
		tlog.Debug("error authenticating", tlog.String("err", err.Error()))
		http.Error(res, "Error authenticating", http.StatusForbidden)
		return
	}

	ret := map[string]string{"name": user.Name}
	body, err := json.Marshal(ret)
	if err != nil {
		tlog.Info("error parsing json", tlog.String("err", err.Error()))
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		return
	}
	r := bytes.NewReader(body)
	_, err = io.Copy(res, r)
	if err != nil {
		tlog.Warn("error writing to response", tlog.String("err", err.Error()))
	}
}
