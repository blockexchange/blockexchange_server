package main

import (
	"blockexchange/db"
	"blockexchange/web"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.Info("Starting")
	db.Init()
	db.Migrate()
	web.Serve()
}
