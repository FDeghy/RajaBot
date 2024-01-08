package bot

import (
	"github.com/AnimeKaizoku/ratelimiter"
	"github.com/FDeghy/RajaGo/raja"
)

var (
	Stations    *raja.Stations
	RateLimiter *ratelimiter.Limiter
)
