package payment

const (
	BaseZibal      = "https://gateway.zibal.ir"
	CreateTokenURL = BaseZibal + "/v1/request"
	BankURL        = BaseZibal + "/start/%v"
	VerifyURL      = BaseZibal + "/v1/verify"
	BadCbMsg       = "خطا در خواندن اطلاعات بازگشتی از درگاه پرداخت."
	OrderNotFound  = "پرداخت یافت نشد." + "\n" +
		"در صورت کسر مبلغ تا 72 ساعت آینده بازگشت میخورد." + "\n" +
		"برای کسب اطلاعات بیشتر با ادمین در تماس باشید."
	OrderFailed = "پرداخت انجام نشد!" + "\n" +
		"در صورت کسر مبلغ تا 72 ساعت آینده بازگشت میخورد." + "\n" +
		"برای کسب اطلاعات بیشتر با ادمین در تماس باشید."
	OrderSuccessful = "پرداخت انجام شد."
	AddSubFailed    = "خطا در اضافه کردن اشتراک!" + "\n" +
		"جهت رفع مشکل با ادمین در ارتباط باشید."
	AddSubSuccessful = "اشتراک شما با موفقیت فعال گردید."
)
