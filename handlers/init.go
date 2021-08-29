package handlers

import (
	"github.com/rancaka/webvitation-api/util"
)

var (
	fb     *util.FirebaseApp
	config util.Config
)

func init() {
	fb = util.GetFirebaseApp()
	config = util.GetConfig()
}
