package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	td "github.com/huangjiahua/tempdesk"
	"github.com/huangjiahua/tempdesk/internal/mock"
	"net/http"
	"testing"
	"time"
)

func TestHMACAuther_AuthUser(t *testing.T) {
	errIsKind := func(t *testing.T, err error, kind string) {
		if err == nil || err.(*AutherError).Kind != kind {
			t.Errorf("Error is not %v: %v", kind, err)
		}
	}

	a := NewHMACAuther()
	us := mock.NewUserService()
	user1 := td.User{Name: "Sam", Key: "password"}

	_ = us.CreateUser(user1)

	req, _ := http.NewRequest(http.MethodGet, "http://example.com/path/", nil)

	_, err := a.AuthUser(req, us)
	errIsKind(t, err, WrongFormat)

	d := time.Now().UTC().Format(http.TimeFormat)
	req.Header.Add("Date", d)
	req.Header.Add("Authorization", "xxx")
	_, err = a.AuthUser(req, us)
	errIsKind(t, err, WrongFormat)

	req.Header.Set("Authorization", "HMAC Tom xxx")
	_, err = a.AuthUser(req, us)
	errIsKind(t, err, NoUser)

	req.Header.Set("Authorization", "HMAC Sam xxx")
	_, err = a.AuthUser(req, us)
	errIsKind(t, err, NotAuthed)

	mac := hmac.New(sha256.New, []byte(user1.Key))
	msg := []byte(req.Method + "\n" + req.URL.Path + "\n" + user1.Name + "\n" + d)
	mac.Write(msg)
	digest := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	req.Header.Set("Authorization", fmt.Sprintf("HMAC Sam %v", digest))

	_, err = a.AuthUser(req, us)
	if err != nil {
		t.Errorf("Should authed here: %v", err)
	}

	req.Header.Set("Date", time.Now().UTC().Add(-10*time.Minute-1*time.Second).Format(http.TimeFormat))
	_, err = a.AuthUser(req, us)
	errIsKind(t, err, Outdated)
}
