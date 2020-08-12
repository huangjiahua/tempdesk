package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	td "github.com/huangjiahua/tempdesk"
	"github.com/huangjiahua/tempdesk/internal/auth"
	thttp "github.com/huangjiahua/tempdesk/internal/http"
	"github.com/huangjiahua/tempdesk/internal/mock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupHMAC(req *http.Request, user *td.User) {
	d := time.Now().UTC().Format(http.TimeFormat)
	mac := hmac.New(sha256.New, []byte(user.Key))
	msg := []byte(req.Method + "\n" + req.URL.Path + "\n" + user.Name + "\n" + d)
	mac.Write(msg)
	digest := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	req.Header.Set("Date", d)
	req.Header.Set("Authorization", fmt.Sprintf("HMAC %v %v", user.Name, digest))
}

func TestUser_ServeHTTP_Get(t *testing.T) {
	h := &User{
		state: &thttp.State{
			Users:  mock.NewUserService(),
			Files:  mock.NewFileService(),
			Auther: auth.NewHMACAuther(),
		},
	}

	ts := httptest.NewServer(h)
	defer ts.Close()

	user1 := td.User{Name: "Sam", Key: "password"}
	_ = h.state.Users.CreateUser(user1)

	res, _ := http.Get(ts.URL + "/")
	if res.StatusCode != http.StatusForbidden {
		t.Fatal("Wrong response")
	}

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/", nil)
	d := time.Now().UTC().Format(http.TimeFormat)
	mac := hmac.New(sha256.New, []byte(user1.Key))
	msg := []byte(req.Method + "\n" + req.URL.Path + "\n" + user1.Name + "\n" + d)
	mac.Write(msg)
	digest := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	req.Header.Set("Date", d)
	req.Header.Set("Authorization", fmt.Sprintf("HMAC Sam %v", digest))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, http.StatusOK, res.StatusCode, "wrong response status") {
		b, _ := ioutil.ReadAll(res.Body)
		_ = res.Body.Close()
		t.Fatalf("Wrong response: %v", string(b))
	}

	b, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()

	m := map[string]string{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Sam", m["name"], "Wrong response")
}

func TestUser_ServeHTTP_POST(t *testing.T) {
	h := &User{
		state: &thttp.State{
			Users:  mock.NewUserService(),
			Files:  mock.NewFileService(),
			Auther: auth.NewHMACAuther(),
		},
	}

	ts := httptest.NewServer(h)
	defer ts.Close()

	info := userSignUpInfo{
		Name: "jack",
		Key:  "key",
		Meta: make(map[string]string),
	}

	body, err := json.Marshal(&info)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(ts.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Wrong response code")

	resp, err = http.Post(ts.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong response status %v", resp.StatusCode)
	}
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Wrong response code")
	msg, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "name already exists\n", string(msg), "Wrong message")
}

func TestUser_ServeHTTP_PUT(t *testing.T) {
	h := &User{
		state: &thttp.State{
			Users:  mock.NewUserService(),
			Files:  mock.NewFileService(),
			Auther: auth.NewHMACAuther(),
		},
	}

	ts := httptest.NewServer(h)
	defer ts.Close()

	user := td.User{
		Name: "jack",
		Key:  "key",
		Meta: make(map[string]string),
	}

	_ = h.state.Users.CreateUser(user)

	info := userUpdateInfo{
		Name: "jack",
		Key:  "new-key",
		Meta: map[string]string{"sex": "male"},
	}

	body, err := json.Marshal(&info)
	if err != nil {
		t.Fatal(err)
	}

	// right update with correct key
	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/", bytes.NewReader(body))
	setupHMAC(req, &user)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, http.StatusOK, res.StatusCode, "wrong response status") {
		b, _ := ioutil.ReadAll(res.Body)
		_ = res.Body.Close()
		t.Fatalf("Wrong response: %v", string(b))
	}

	user1, _ := h.state.Users.User("jack")
	assert.Equal(t, "new-key", user1.Key, "wrong key")

	// update with outdated key
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusForbidden, res.StatusCode, "wrong response status")

	user.Key = "new-key"
	req, _ = http.NewRequest(http.MethodPut, ts.URL+"/", bytes.NewReader(body))
	setupHMAC(req, &user)

	// update with renewed key
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, http.StatusOK, res.StatusCode, "wrong response status") {
		b, _ := ioutil.ReadAll(res.Body)
		_ = res.Body.Close()
		t.Fatalf("Wrong response: %v", string(b))
	}
}

func TestUser_ServeHTTP_Delete(t *testing.T) {
	h := &User{
		state: &thttp.State{
			Users:  mock.NewUserService(),
			Files:  mock.NewFileService(),
			Auther: auth.NewHMACAuther(),
		},
	}

	ts := httptest.NewServer(h)
	defer ts.Close()

	user := td.User{
		Name: "jack",
		Key:  "key",
		Meta: make(map[string]string),
	}

	_ = h.state.Users.CreateUser(user)

	info := userUpdateInfo{
		Name: "jack",
		Key:  "",
	}

	body, err := json.Marshal(&info)
	if err != nil {
		t.Fatal(err)
	}

	// right update with correct key
	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/", bytes.NewReader(body))
	setupHMAC(req, &user)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, http.StatusOK, res.StatusCode, "wrong response status") {
		b, _ := ioutil.ReadAll(res.Body)
		_ = res.Body.Close()
		t.Fatalf("Wrong response: %v", string(b))
	}

	_, fnd := h.state.Users.User("jack")
	assert.Equal(t, false, fnd, "does not delete user")

	// redo delete
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusForbidden, res.StatusCode, "wrong response status")
}
