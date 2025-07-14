package methodutil

import "ems/config"

func AccessUuidCacheKey(accessUuid string) string {
	return config.Redis().MandatoryPrefix + config.Redis().AccessUuidPrefix + accessUuid
}

func RefreshUuidCacheKey(refreshUuid string) string {
	return config.Redis().MandatoryPrefix + config.Redis().RefreshUuidPrefix + refreshUuid
}
