package main

import (
	"log/slog"
	"time"

	"github.com/BANKA2017/bing-daily/dbio"
	"github.com/BANKA2017/bing-daily/dbio/model"
	"github.com/BANKA2017/bing-daily/server"
	"github.com/BANKA2017/bing-daily/workers"
	"github.com/kdnetwork/code-snippet/go/log"
	"gorm.io/gorm/logger"
)

func main() {
	dbio.InitLogger()

	dbio.InitEnv()

	logLevel := logger.Error
	if dbio.TestMode {
		logLevel = logger.Info
		log.SlogLevel.Set(slog.LevelDebug)
	}

	dbio.GormDB.LogLevel = logLevel
	dbio.GormDB.ServicePrefix = "bing-daily"
	dbio.GormDB.SetLogger(logger.NewSlogLogger(
		slog.Default(), logger.Config{
			// copy from gorm default logger config
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  dbio.GormDB.LogLevel,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		},
	))

	if err := dbio.GormDB.SetDBAuth(dbio.DBUser, dbio.DBPassword, dbio.DBHost, dbio.DBDatabase, dbio.DBCert).Connect(); err != nil {
		log.Fatal("connect to db failed", "error", err)
	}

	dbio.GormMemCacheDB.LogLevel = logLevel
	dbio.GormMemCacheDB.ServicePrefix = "bing-daily-cache"
	if err := dbio.GormMemCacheDB.SetDBPath("file::memory:?cache=shared").Connect(); err != nil {
		log.Fatal("connect to cache db failed", "error", err)
	}

	if err := dbio.GormMemCacheDB.W.AutoMigrate(&model.Img2{}, &model.Color{}); err != nil {
		log.Fatal("init cache db failed", "error", err)
	}

	// init cache
	var img2 []*model.Img2
	var color2 []*model.Color

	dbio.GormDB.R.Model(&model.Img2{}).Find(&img2)

	for i := 0; i < len(img2); i += 300 {
		end := min(i+300, len(img2))
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
	if err := workers.GetImgsWorker(dbio.B2ApplicationKeyId, dbio.B2ApplicationKey, dbio.WorkersLocale); err != nil {
		slog.Error("get images worker failed", "error", err)
	}

	if dbio.Addr != "" {
		go server.Server(dbio.Addr)
	}
	for range oneHourInterval.C {
		if err := workers.GetImgsWorker(dbio.B2ApplicationKeyId, dbio.B2ApplicationKey, dbio.WorkersLocale); err != nil {
			slog.Error("get images worker failed", "error", err)
		}

		workers.UpdateLatestDate()
	}
}
