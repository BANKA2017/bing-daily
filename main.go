package main

import (
	"log"
	"time"

	"github.com/BANKA2017/bing-daily/dbio"
	"github.com/BANKA2017/bing-daily/dbio/model"
	"github.com/BANKA2017/bing-daily/server"
	"github.com/BANKA2017/bing-daily/workers"
	"gorm.io/gorm/logger"
)

// server
//var host string

var err error

var testmode bool

func main() {
	dbio.InitEnv()

	logLevel := logger.Error
	if testmode {
		logLevel = logger.Info
	}

	dbio.GormDB.R, dbio.GormDB.W, err = dbio.ConnectToMySQL(dbio.DBUser, dbio.DBPassword, dbio.DBHost, dbio.DBDatabase, dbio.DBCert, logLevel, "bing-daily")
	if err != nil {
		log.Fatal(err)
	}

	dbio.GormMemCacheDB.R, dbio.GormMemCacheDB.W, err = dbio.ConnectToSQLite("file::memory:?cache=shared", logLevel, "bing-daily")
	if err != nil {
		log.Fatal(err)
	}

	if err := dbio.GormMemCacheDB.W.AutoMigrate(&model.Img2{}, &model.Color{}); err != nil {
		log.Fatal(err)
	}

	// init cache
	var img2 []*model.Img2
	var color2 []*model.Color

	dbio.GormDB.R.Model(&model.Img2{}).Find(&img2)

	for i := 0; i < len(img2); i += 300 {
		end := i + 300
		if end > len(img2) {
			end = len(img2)
		}
		dbio.GormMemCacheDB.W.Create(img2[i:end])
	}

	dbio.GormDB.R.Model(&model.Color{}).Find(&color2)
	for i := 0; i < len(color2); i += 300 {
		end := i + 300
		if end > len(color2) {
			end = len(color2)
		}
		dbio.GormMemCacheDB.W.Create(color2[i:end])
	}

	oneHourInterval := time.NewTicker(time.Hour)
	defer oneHourInterval.Stop()

	workers.UpdateLatestDate()
	go workers.GetImgsWorker(dbio.B2ApplicationKeyId, dbio.B2ApplicationKey, dbio.WorkersLocale)

	if dbio.Addr != "" {
		go server.Server(dbio.Addr)
	}
	for {
		select {
		case <-oneHourInterval.C:
			if err := workers.GetImgsWorker(dbio.B2ApplicationKeyId, dbio.B2ApplicationKey, dbio.WorkersLocale); err != nil {
				log.Println(err)
			}

			workers.UpdateLatestDate()
		}
	}
}
