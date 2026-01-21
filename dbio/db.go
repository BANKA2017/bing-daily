package dbio

import (
	"github.com/kdnetwork/code-snippet/go/db"
)

var GormDB = new(db.GormDBCtx).SetDBMode(db.DBModeMySQL)
var GormMemCacheDB = (&db.GormDBCtx{AllowMemoryMode: true}).SetDBMode(db.DBModeSQLite)
