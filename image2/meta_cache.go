package image2

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var ImageMetaCache = ttlcache.New(
	ttlcache.WithCapacity[string, ImageMeta](200),
	ttlcache.WithTTL[string, ImageMeta](time.Minute*10),
)
