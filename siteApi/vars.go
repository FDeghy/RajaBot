package siteapi

import "errors"

const (
	BASE_URL = "https://ticket.rai.ir"
)

var (
	ErrJsonDecode  = errors.New("cannot decode json")
	ErrBadHttpCode = errors.New("http response is not 200")
)
