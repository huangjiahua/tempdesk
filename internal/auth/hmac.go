package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	td "github.com/huangjiahua/tempdesk"
	"net/http"
	"strings"
	"time"
)

type HMACAuther struct {
}

func NewHMACAuther() HMACAuther {
	return HMACAuther{}
}

func (j HMACAuther) AuthUser(req *http.Request, us td.UserService) (td.User, error) {
	a := req.Header.Get("Authorization")
	d := req.Header.Get("Date")
	if len(a) == 0 || len(d) == 0 {
		return td.User{}, &AutherError{WrongFormat, "Missing Header"}
	}

	t, err := http.ParseTime(d)
	if err != nil {
		return td.User{}, &AutherError{WrongFormat, "Wrong Date Header Format"}
	}

	now := time.Now().UTC()
	if now.Sub(t) > 10*time.Minute || t.Sub(now) > 10*time.Minute {
		return td.User{}, &AutherError{Outdated, "Message Is Not Valid"}
	}

	fields := strings.Fields(a)
	if len(fields) != 3 || fields[0] != "HMAC" {
		return td.User{}, &AutherError{WrongFormat, "Wrong Authorization Header Format"}
	}

	username := fields[1]

	user, ok := us.User(username)
	if !ok {
		return user, &AutherError{NoUser, "Cannot Find User"}
	}

	msg := req.Method + "\n" + req.URL.Path + "\n" + username + "\n" + d
	if !ValidDigest([]byte(msg), fields[2], []byte(user.Key)) {
		return td.User{}, &AutherError{NotAuthed, "Not Authed"}
	}

	return user, nil
}

func ValidDigest(message []byte, digest string, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return expected == digest
}
