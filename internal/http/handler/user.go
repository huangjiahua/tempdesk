package handler

import (
	"bytes"
	"encoding/json"
	td "github.com/huangjiahua/tempdesk"
	thttp "github.com/huangjiahua/tempdesk/internal/http"
	tlog "github.com/huangjiahua/tempdesk/pkg/log"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	ErrorInternalServer = "internal server error"
	ErrorParsingJson    = "error parsing json"
	ErrorEncodingJson   = "error encoding json"
	ErrorParsingBody    = "error parsing body"
	ErrorAuthenticating = "error authenticating"
	ErrorWritingResp    = "error writing to response"
	ErrorCloseReader    = "error close reader"
	ErrorCreatingUser   = "error creating user"
	ErrorUpdatingUser   = "error updating user"
	ErrorDeletingUser   = "error deleting user"

	ActionUpdate = "update"
	ActionDelete = "delete"
)

type User struct {
	state *thttp.State
}

func (u *User) ServeGetUser(res http.ResponseWriter, req *http.Request) {
	user, err := u.state.AuthUser(req)
	if err != nil {
		tlog.Debug(ErrorAuthenticating, tlog.Err(err))
		http.Error(res, ErrorAuthenticating, http.StatusForbidden)
		return
	}

	ret := map[string]string{"name": user.Name}
	body, err := json.Marshal(ret)
	if err != nil {
		tlog.Info(ErrorEncodingJson, tlog.Err(err))
		http.Error(res, ErrorInternalServer, http.StatusInternalServerError)
		return
	}
	r := bytes.NewReader(body)
	_, err = io.Copy(res, r)
	if err != nil {
		tlog.Warn(ErrorWritingResp, tlog.Err(err))
	}
}

func (u *User) ServeAddUser(res http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		tlog.Debug(ErrorParsingBody, tlog.Err(err))
		http.Error(res, ErrorParsingBody, http.StatusBadRequest)
		return
	}
	if err = req.Body.Close(); err != nil {
		tlog.Debug(ErrorCloseReader, tlog.Err(err))
	}

	info, err := parseUserSignUpInfo(body)
	if err != nil {
		tlog.Debug(ErrorParsingJson, tlog.Err(err))
		http.Error(res, ErrorParsingJson, http.StatusBadRequest)
		return
	}
	user := td.User{
		Name: info.Name,
		Key:  info.Key,
		Meta: info.Meta,
	}

	err = u.state.Users.CreateUser(user)
	if err != nil {
		tlog.Debug(ErrorCreatingUser, tlog.Err(err))
		switch err.(*td.UserServiceError).Kind {
		case td.NameAlreadyExists:
			http.Error(res, err.Error(), http.StatusBadRequest)
		default:
			http.Error(res, ErrorCreatingUser, http.StatusBadRequest)
		}
		return
	}

	tlog.Info("add new user",
		tlog.String("name", user.Name))

	res.WriteHeader(http.StatusOK)
}

func (u *User) ServeUpdateUser(res http.ResponseWriter, req *http.Request, action string) {
	user, err := u.state.AuthUser(req)
	if err != nil {
		tlog.Debug(ErrorAuthenticating, tlog.Err(err))
		http.Error(res, ErrorAuthenticating, http.StatusForbidden)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		tlog.Debug(ErrorParsingBody, tlog.Err(err))
		http.Error(res, ErrorParsingBody, http.StatusBadRequest)
		return
	}

	info, err := parseUpdateInfo(body)
	if err != nil {
		tlog.Debug(ErrorParsingJson, tlog.Err(err))
		http.Error(res, ErrorParsingJson, http.StatusBadRequest)
		return
	}

	upd := td.User{
		Name: info.Name,
		Key:  info.Key,
		Meta: info.Meta,
	}

	if action == ActionUpdate {
		err = u.state.Users.UpdateUser(upd)
		if err != nil {
			tlog.Debug(ErrorUpdatingUser, tlog.Err(err))
			http.Error(res, ErrorUpdatingUser, http.StatusBadRequest)
			return
		}
		tlog.Info("update user request",
			tlog.String("exe", user.Name),
			tlog.String("target", info.Name))
	} else {
		// delete
		err = u.state.Users.DeleteUser(upd)
		if err != nil {
			tlog.Debug(ErrorDeletingUser, tlog.Err(err))
			http.Error(res, ErrorDeletingUser, http.StatusBadRequest)
			return
		}
		tlog.Info("delete user request",
			tlog.String("exe", user.Name),
			tlog.String("target", info.Name))
	}

	res.WriteHeader(http.StatusOK)
}

func (u *User) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		tlog.Debug("error parsing url", tlog.Err(err))
		http.Error(res, "Error parsing url arguments or form", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		u.ServeGetUser(res, req)
	case http.MethodPost:
		u.ServeAddUser(res, req)
	case http.MethodPut:
		u.ServeUpdateUser(res, req, ActionUpdate)
	case http.MethodDelete:
		u.ServeUpdateUser(res, req, ActionDelete)
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

type userUpdateInfo struct {
	Name string            `json:"name"`
	Key  string            `json:"password,omitempty"`
	Meta map[string]string `json:"meta,omitempty"`
}

func parseUpdateInfo(data []byte) (userUpdateInfo, error) {
	var i userUpdateInfo
	err := json.Unmarshal(data, &i)
	return i, err
}
