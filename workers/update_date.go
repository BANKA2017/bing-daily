package workers

import (
	"github.com/BANKA2017/bing-daily/bing"
	"github.com/BANKA2017/bing-daily/dbio"
	"github.com/BANKA2017/bing-daily/dbio/model"
)

func UpdateLatestDate() error {
	var DateList []*model.Img2

	if err := dbio.GormDB.R.Model(&model.Img2{}).Select("MAX(date) AS date, market").Group("market").Scan(&DateList).Error; err != nil {
		return err
	}

	for _, d := range DateList {
		bing.LatestDate[d.Market] = int64(d.Date)
	}

	return nil
}
