package main

import (
	"blockexchange/db"
	"blockexchange/web"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)
	logrus.Info("Starting")
	db.Init()
	db.Migrate()
	web.Serve()
}
