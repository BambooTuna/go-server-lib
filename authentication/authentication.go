package authentication

import (
	"github.com/BambooTuna/go-server-lib/session"
	"github.com/google/uuid"
	"time"
)

type ActivatorUseCase struct {
	Activator      Activator
	ActivateMailer ActivateMailer
}

func (a ActivatorUseCase) IssueCode(id string, mailAddress string) (string, error) {
	if code, err := a.Activator.IssueCode(id); err != nil {
		return "", err
	} else if err := a.ActivateMailer.SendActivateCode(code, mailAddress); err != nil {
		return "", err
	} else {
		return code, nil
	}
}

func (a ActivatorUseCase) Activate(code string) (string, error) {
	return a.Activator.Activate(code)
}

func DefaultActivatorUseCase(sessionStorageDao session.SessionStorageDao, activateMailer ActivateMailer) ActivatorUseCase {
	codeGenerator := func() (string, error) {
		uuidObj, err := uuid.NewUUID()
		data := []byte("wnw8olzvmjp0x6j7ur8vafs4jltjabi0")
		uuidObj2 := uuid.NewMD5(uuidObj, data)
		return uuidObj2.String(), err
	}

	return ActivatorUseCase{
		Activator: Activator{
			SessionStorageDao: sessionStorageDao,
			ExpirationDate:    time.Duration(1) * time.Hour,
			CodeGenerator:     codeGenerator,
		},
		ActivateMailer: activateMailer,
	}
}
