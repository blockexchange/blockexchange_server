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
	db_, err := db.Init()
	if err != nil {
		panic(err)
	}
	db.Migrate(db_.DB)
	web.Serve(db_)
}
