package bot

import (
	"github.com/ALiwoto/ratelimiter"
	"github.com/FDeghy/RajaGo/raja"
)

var (
	Stations    *raja.Stations
	RateLimiter *ratelimiter.Limiter
)
