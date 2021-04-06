package main

import (
	"github.com/sirupsen/logrus"

	"github.com/BenJetson/CPSC491-project/go/app/api"
	"github.com/BenJetson/CPSC491-project/go/app/db"
	"github.com/BenJetson/CPSC491-project/go/app/etsy"
)

func main() {
	logger := logrus.New()

	dbCfg, err := db.NewConfigFromEnv()
	if err != nil {
		logger.Fatalln(err)
	}

	db, err := db.NewDataStore(logger, dbCfg)
	if err != nil {
		logger.Fatalln(err)
	}

	svrCfg, err := api.NewConfigFromEnv()
	if err != nil {
		logger.Fatalln(err)
	}

	cv, err := etsy.NewClientFromEnv()
	if err != nil {
		logger.Fatalln(err)
	}

	svr, err := api.NewServer(logger, db, cv, svrCfg)
	if err != nil {
		logger.Fatalln(err)
	}

	err = svr.Start()
	if err != nil {
		logger.Fatalln(err)
	}
}
