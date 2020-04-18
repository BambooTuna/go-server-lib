package authentication

import (
	"github.com/BambooTuna/go-server-lib/session"
	"time"
)

type Activator struct {
	SessionStorageDao session.SessionStorageDao
	ExpirationDate    time.Duration
	CodeGenerator     func() (string, error)
}

// 引数のIdに対するアクティベート用コードが発行されます
func (a Activator) IssueCode(id string) (string, error) {
	if code, err := a.CodeGenerator(); err != nil {
		return "", err
	} else if err := a.SessionStorageDao.Store(code, id, a.ExpirationDate); err != nil {
		return "", err
	} else {
		return code, nil
	}
}

// アクティベート用コードが有効であれば元のIdが返される
func (a Activator) Activate(code string) (string, error) {
	if id, err := a.SessionStorageDao.Find(code); err != nil {
		return "", err
	} else {
		a.SessionStorageDao.Remove(code)
		return id, nil
	}
}
