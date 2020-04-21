package session_v2

import "time"

type SessionStorageDao interface {
	Store(key, value string, expiration time.Duration) (string, error)
	Find(key string) (string, error)
	Remove(key string) (int64, error)
	Refresh(key string, expiration time.Duration) (bool, error)
}
