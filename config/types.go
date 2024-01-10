package config

type Config struct {
	Bot struct {
		Token                  string   `toml:"token"`
		HttpURI                string   `toml:"http_proxy"`
		Timeout                int      `toml:"timeout"`
		StationsButtonsPerPage int      `toml:"stations_buttons_per_page"`
		FavoriteStations       []string `toml:"favorite_stations"`
		Admins                 []int64  `toml:"admins"`
		UserLimit              int      `toml:"user_limit"`
		VipLimit               int      `toml:"vip_limit"`
	} `toml:"BOT"`
	Database struct {
		Name string `toml:"name"`
	} `toml:"DB"`
	Raja struct {
		Timeout    int   `toml:"timeout"`
		CheckEvery int   `toml:"check_every"`
		AlertEvery int64 `toml:"alert_every"`
		Buffer     int   `toml:"buffer"`
		Worker     int   `toml:"worker"`
	} `toml:"RAJA"`
	Prometheus struct {
		AddressBind string `toml:"address_bind"`
	} `toml:"PROMETHEUS"`
}
