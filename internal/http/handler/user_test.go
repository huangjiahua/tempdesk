package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	td "github.com/huangjiahua/tempdesk"
	"github.com/huangjiahua/tempdesk/internal/auth"
	thttp "github.com/huangjiahua/tempdesk/internal/http"
	"github.com/huangjiahua/tempdesk/internal/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUser_ServeHTTP(t *testing.T) {
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

	if res.StatusCode != http.StatusOK {
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

	if m["name"] != "Sam" {
		t.Error("Wrong response")
	}
}
