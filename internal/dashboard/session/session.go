package session

import (
	"crypto/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

func init() {
	_, err := rand.Read(jwtkey)
	if err != nil {
		panic(err)
	}
}

var (
	sessions = make(map[string]*ClientInfo, 0)
	jwtkey   = make([]byte, 128)
)

const (
	SESSION_EXP_TIME = time.Minute * 5
)

type ClientInfo struct {
	IPAddr string    `json:"ipAddr"`
	Root   bool      `json:"root"`
	Iat    time.Time `json:"iat "`
	exp    time.Time
}

func NewSession(ipAddr string) (tokenStr string, err error) {
	sid, err := generateSID()
	if err != nil {
		return "", err
	}

	exp := time.Now().Add(SESSION_EXP_TIME)
	sessions[tokenStr] = &ClientInfo{
		IPAddr: ipAddr,
		Iat:    time.Now(),
		exp:    exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid": string(sid),
	})
	tokenStr, err = token.SignedString(jwtkey)
	if err != nil {
		return "", err
	}

	// clear expired sessions every time new session is created
	go clearExpiredSessions()
	return tokenStr, nil
}

func generateSID() ([]byte, error) {
	sid := make([]byte, 32)
	for {
		_, err := rand.Read(sid)
		if err != nil {
			return nil, err
		}
		if _, ok := sessions[string(sid)]; !ok {
			break
		}
	}
	return sid, nil
}

func clearExpiredSessions() {
	for sid, cinfo := range sessions {
		if cinfo.exp.Before(time.Now()) {
			delete(sessions, sid)
		}
	}
}

func GetAndExtendSession(tokenStr string) (*ClientInfo, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		return nil, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}

	sid, ok := claims["sid"].(string)
	if !ok {
		return nil, false
	}

	cinfo, ok := sessions[sid]
	if cinfo.exp.Before(time.Now()) {
		return nil, false
	}
	cinfo.exp = time.Now().Add(SESSION_EXP_TIME)
	return cinfo, ok
}
