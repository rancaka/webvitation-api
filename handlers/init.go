package handlers

import (
	"github.com/rancaka/webvitation-web/util"
)

var (
	fb *util.FirebaseApp
)

func init() {
	fb = util.GetFirebaseApp()
}
