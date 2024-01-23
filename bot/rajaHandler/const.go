package rajaHandler

const (
	StartMsg = "خوش آمدید!" + "\n" +
		"ساعت های قطعی رجا:" + "\n" +
		"00:00 -> 00:15" + "\n" +
		"06:30 -> 06:45" + "\n" +
		"13:30 -> 13:45" + "\n" +
		"19:30 -> 19:45"
	SrcMsg           = "لطفا ایستگاه مبدا را انتخاب کنید"
	DstMsg           = "لطفا ایستگاه مقصد را انتخاب کنید"
	StationsLoadErr  = "خطا در لود کردن ایستگاه ها"
	TrainListLoadErr = "خطا در لود کردن قطار ها"
	NextPage         = "صفحه بعد ⬅️"
	PreviousPage     = "➡️ صفحه قبل"
	NextMonth        = "ماه بعد ⬅️"
	PreviousMonth    = "➡️ ماه قبل"
	PageN            = "صفحه %d"
	FavSign          = " " + "⭐️"
	AnError          = "خطای ناشناخته با ادمین تماس بگیرید"
	DayMsg           = "تاریخ را انتخاب کنید" + "\n" + "MMM yyyy"
	TimeFormat       = "E dd MMM yyyy"
	TrainSelMsg      = TimeFormat
	StateErr         = "خطای مرحله! لطفا فعالیت قبلی را کامل کنید یا آن را کنسل کنید"
	TrainButtonText  = "%s - %s - %s تومان"
	GetTrainsInfoMsg = "در حال دریافت لیست قطار ها"
	CancelMsg        = "دستور کنسل را فرستاده و مجدد تلاش کنید."
	successfulCreate = "درخواست با موفقیت ثبت شد"
	LimitReached     = "شما اجازه ثبت درخواست فعال بیشتر از این را ندارید. در صورت تمایل به افزایش محدودیت با ادمین تماس بگیرید."
	RajaErr          = "خطا در رجا"
	RajaErr2         = "متأسفیم! با توجه به قطعی سرور راه آهن، پاسخی برای درخواست شما یافت نشد."
	OldTrErr         = "قطار حرکت کرده است."
	OldDateErr       = "تاریخ قدیمی است."
	CancelOkMsg      = "عملیات ناقص فعلی کنسل شد."
	NilButton        = "دکمه پوچ رو زدی!"
	TrainNotFound    = "قطار پیدا نشد! (ممکن است باگ رجا باشد، دقایقی دیگر تست کنید)"
	EmptyTrainWR     = "درخواست فعالی ندارید\\."
	ListReqs         = "لیست درخواست ها:"
	CancButtonTxt    = "❌ حذف درخواست (%v)"
	CancOkAlert      = "با موفقیت حذف شد"
	NoHaveSub        = "اشتراک فعال ندارید" + "\n" +
		"برای مشاهده وضعیت اشتراک دستور /sub را ارسال کنید"
)
