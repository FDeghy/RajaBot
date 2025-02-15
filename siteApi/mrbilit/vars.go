package mrbilit

import "errors"

const (
	BASE_URL = "https://train.mrbilit.com"
)

var (
	ErrJsonDecode  = errors.New("cannot decode json")
	ErrBadHttpCode = errors.New("http response is not 200")
)
