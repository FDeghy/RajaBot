package payment

import "errors"

var (
	ErrBadCode                     = errors.New("bad code from nextpay")
	ErrUncompletedTransactionFound = errors.New("the user have a uncompleted trasnaction")
)
