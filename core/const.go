package core

const (
	AlertMsg          = "قطار %s روز %s از %s به %s با ظرفیت %d نفر خالی شد!"
	TrainDate         = "E dd MMM yyyy"
	ExpireMsg         = "درخواست قطار ساعت %s روز %s از %s به %s منقضی شد."
	CancelMsg         = "درخواست قطار ساعت %s روز %s از %s به %s کنسل شد."
	RajaSearchDateFmt = "yyyyMMdd"
	RajaSearchURL     = "https://www.raja.ir/search?adult=1&child=0&infant=0&movetype=1&ischarter=false&fs=%v&ts=%v&godate=%v&tickettype=Family&returndate=&numberpassenger=1&mode=Train"
	RajaSearchButTxt  = "باز کردن درخواست در سایت رجا"

	SubscriptionStatusMsg = "وضعیت اشتراک شما: %v" + "\n" +
		"تاریخ عضویت: %v" + "\n" +
		"تاریخ انقضا: %v" + "\n" +
		"محدودیت درخواست های همزمان: %v" + "\n"
	Enabled    = "فعال"
	Disabled   = "غیرفعال"
	Unkown     = "نامشخص"
	TimeFormat = "dd MMM yyyy"
)
