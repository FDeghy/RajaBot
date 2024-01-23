package rajaHandler

import (
	"RajaBot/config"
	"strings"
)

func favoriteSign(st string) string {
	for _, i := range config.Cfg.Bot.FavoriteStations {
		if strings.EqualFold(st, i) {
			return FavSign
		}
	}
	return ""
}
