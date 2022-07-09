package main

import (
	martini_init "GolangTD/martini.init"
	"GolangTD/session"

	_ "github.com/mattn/go-sqlite3"
)

var SessionInMemory *session.Session

func main() {
	session.SessionInMemory = session.NewSession()
	martini_init.CommonHandler(martini_init.InitMartini())
}
