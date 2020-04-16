package session

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type DefaultSession struct {
	Dao      SessionStorageDao
	Settings SessionSettings
}

func GenerateUUID() (string, error) {
	uuidObj, err := uuid.NewUUID()
	data := []byte("wnw8olzvmjp0x6j7ur8vafs4jltjabi0")
	uuidObj2 := uuid.NewMD5(uuidObj, data)
	return uuidObj2.String(), err
}

func (s DefaultSession) SetSession(f func(*gin.Context) *string) func(*gin.Context) {
	return func(c *gin.Context) {
		plainToken := f(c)
		if plainToken == nil {
			return
		}
		id, err := GenerateUUID()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &jwt.StandardClaims{Id: id})
		SignedToken, err := token.SignedString([]byte(s.Settings.Secret))
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := s.Dao.Store(id, *plainToken, s.Settings.ExpirationDate); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Header(s.Settings.SetAuthHeaderName, SignedToken)
	}
}

func (s DefaultSession) RequiredSession(f func(*gin.Context, string)) func(*gin.Context) {
	return func(c *gin.Context) {
		claimsId := s.ParseClaimsId(c)
		if claimsId == "" {
			return
		}
		token, err := s.Dao.Find(claimsId)
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		f(c, *token)
	}
}

func (s DefaultSession) OptionalRequiredSession(f func(*gin.Context, *string)) func(*gin.Context) {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(s.Settings.AuthHeaderName)
		if tokenString == "" {
			f(c, nil)
			return
		}
		var claims jwt.StandardClaims
		_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Settings.Secret), nil
		})
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		token, err := s.Dao.Find(claims.Id)
		if err != nil {
			c.Status(http.StatusForbidden)
			return
		}
		f(c, token)
	}
}

func (s DefaultSession) InvalidateSession(f func(*gin.Context)) func(*gin.Context) {
	return s.RequiredSession(func(c *gin.Context, i string) {
		claimsId := s.ParseClaimsId(c)
		if claimsId == "" {
			return
		}

		if err := s.Dao.Remove(claimsId); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		f(c)
	})
}

func (s DefaultSession) ParseClaimsId(c *gin.Context) string {
	tokenString := c.GetHeader(s.Settings.AuthHeaderName)
	if tokenString == "" {
		c.Status(http.StatusForbidden)
		return ""
	}
	var claims jwt.StandardClaims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Settings.Secret), nil
	})
	if err != nil {
		c.Status(http.StatusForbidden)
		return ""
	}
	return claims.Id
}
