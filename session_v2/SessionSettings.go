package session_v2

import (
	"github.com/google/uuid"
	"time"
)

type SessionSettings struct {
	Secret            string
	SetAuthHeaderName string
	AuthHeaderName    string
	ExpirationDate    time.Duration
	IDGenerator       func() (string, error)
}

func DefaultSessionSettings(secret string) SessionSettings {
	return SessionSettings{
		Secret:            secret,
		SetAuthHeaderName: "set-authorization",
		AuthHeaderName:    "authorization",
		ExpirationDate:    time.Duration(1) * time.Hour,
		IDGenerator: func() (string, error) {
			uuidObj, err := uuid.NewUUID()
			data := []byte("wnw8olzvmjp0x6j7ur8vafs4jltjabi0")
			uuidObj2 := uuid.NewMD5(uuidObj, data)
			return uuidObj2.String(), err
		},
	}
}
