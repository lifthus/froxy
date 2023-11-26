package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
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
	ssmu     = sync.Mutex{}
	sessions = make(map[string]*ClientInfo, 0)
	jwtkey   = make([]byte, 128)
)

const (
	SESSION_EXP_TIME = time.Minute * 5
)

// for embedding in request context
type cinfokey string

const Cinfokey cinfokey = "cinfokey"

type ClientInfo struct {
	IPAddr string    `json:"ipAddr"`
	Root   bool      `json:"root"`
	Iat    time.Time `json:"iat "`
	exp    time.Time
}

func NewSession(ipAddr string) (tokenStr string, cinfo *ClientInfo, err error) {
	sidb, err := generateSID()
	if err != nil {
		return "", nil, err
	}

	sid := hex.EncodeToString(sidb)

	exp := time.Now().Add(SESSION_EXP_TIME)
	newCinfo := &ClientInfo{
		IPAddr: ipAddr,
		Iat:    time.Now(),
		exp:    exp,
	}

	ssmu.Lock()
	sessions[string(sid)] = newCinfo
	ssmu.Unlock()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid": string(sid),
	})
	tokenStr, err = token.SignedString(jwtkey)
	if err != nil {
		return "", nil, err
	}

	// clear expired sessions every time new session is created
	go clearExpiredSessions()
	return tokenStr, newCinfo, nil
}

func generateSID() ([]byte, error) {
	sid := make([]byte, 4)
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
			ssmu.Lock()
			delete(sessions, sid)
			ssmu.Unlock()
		}
	}
}

func GetAndExtendSession(tokenStr string) (*ClientInfo, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	sid, ok := claims["sid"].(string)
	if !ok {
		return nil, fmt.Errorf("no sid in token")
	}

	cinfo, ok := sessions[sid]
	if !ok || cinfo.exp.Before(time.Now()) {
		return nil, fmt.Errorf("invalid sid")
	}

	cinfo.exp = time.Now().Add(SESSION_EXP_TIME)
	return cinfo, nil
}
