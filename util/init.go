package util

import (
	"github.com/jmoiron/sqlx"
)

var (
	GetConfig      func() Config
	GetDB          func() *sqlx.DB
	GetFirebaseApp func() *FirebaseApp
)

func init() {
	GetConfig = getConfig
	GetDB = getDB
	GetFirebaseApp = getFirebaseApp
}
