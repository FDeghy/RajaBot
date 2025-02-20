package tlog

type Mode uint8

const (
	Payment Mode = iota
	NewTrain
)

const (
	PaymentMsg = "پرداخت جدید" + "\n" +
		"کاربر: %v" + "\n" +
		"توضیح: %v" + "\n" +
		"تاریخ: %v"
	NewTrainMsg = "درخواست قطار جدید" + "\n" +
		"کاربر: %v" + "\n" +
		"قطار: %v" + "\n" +
		"تاریخ: %v"
)
