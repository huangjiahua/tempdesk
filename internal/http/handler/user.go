package handler

import (
	"bytes"
	"encoding/json"
	td "github.com/huangjiahua/tempdesk"
	thttp "github.com/huangjiahua/tempdesk/internal/http"
	tlog "github.com/huangjiahua/tempdesk/internal/log"
	"io"
	"io/ioutil"
	"net/http"
)

type User struct {
	state *thttp.State
}

func (u *User) ServeGetUser(res http.ResponseWriter, req *http.Request) {
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

func (u *User) ServeAddUser(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		tlog.Debug("error parsing body", tlog.String("err", err.Error()))
		http.Error(res, "error parsing body", http.StatusBadRequest)
		return
	}
	info, err := parseUserSignUpInfo(body)
	if err != nil {
		tlog.Debug("error parsing submitted json", tlog.String("err", err.Error()))
		http.Error(res, "error parsing submitted json", http.StatusBadRequest)
		return
	}
	user := td.User{
		Name: info.Name,
		Key:  info.Key,
		Meta: info.Meta,
	}

	err = u.state.Users.CreateUser(user)
	if err != nil {
		// TODO: handle each kind of error
		tlog.Debug("error creating new user", tlog.String("err", err.Error()))
		http.Error(res, "error creating new user", http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (u *User) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		tlog.Debug("error parsing url", tlog.String("err", err.Error()))
		http.Error(res, "Error parsing url arguments or form", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		u.ServeGetUser(res, req)
	case http.MethodPost:
		u.ServeAddUser(res, req)
	default:
		tlog.Debug("unsupported method", tlog.String("method", req.Method))
		http.Error(res, "method not supported", http.StatusMethodNotAllowed)
	}
}

type userSignUpInfo struct {
	Name string            `json:"name"`
	Key  string            `json:"password"`
	Meta map[string]string `json:"meta,omitempty"`
}

func parseUserSignUpInfo(data []byte) (userSignUpInfo, error) {
	var i userSignUpInfo
	err := json.Unmarshal(data, &i)
	return i, err
}
