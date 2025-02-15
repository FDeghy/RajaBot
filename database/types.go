package database

type TgUser struct {
	UserID int64  `json:"user_id" gorm:"primaryKey"`
	IsVip  bool   `json:"is_vip"`
	State  string `json:"state"`
}

type TrainWR struct {
	Id      uint64 `json:"id" gorm:"primaryKey"`
	UserID  int64  `json:"user_id"`
	Day     int64  `json:"day"`
	TrainId int    `json:"train_id"`
	Hour    string `json:"hour"`
	Src     int    `json:"src"` // in homei mode is route id
	Dst     int    `json:"dst"` // -1 mean homei
	IsDone  bool   `json:"is_done"`
	ThrdApp int    `json:"thrd_app"`
}

type Subscription struct {
	UserID         int64 `json:"user_id" gorm:"primaryKey"`
	IsTrial        bool  `json:"is_trial"`
	ExpirationDate int64 `json:"expiration_date"`
	RegisteryDate  int64 `json:"registery_date"`
	IsEnabled      bool  `json:"is_enabled"`
}

type Payment struct {
	Id            uint64 `json:"id" gorm:"primaryKey"`
	UserID        int64  `json:"user_id"`
	OrderID       string `json:"order_id"`
	Price         uint   `json:"price"`
	TransID       string `json:"trans_id"`
	CardNumber    string `json:"card_number"`
	CustomerPhone string `json:"customer_phone"`
	ShaparakRefId string `json:"shaparak_ref_id"`
	OrderDate     string `json:"order_date"`
	StatusCode    int    `json:"status_code"`
	IsDone        bool   `json:"is_done"`
}

type RTTrain struct {
	Id     uint64 `json:"id" gorm:"primaryKey"`
	Src    string `json:"src"`
	Dst    string `json:"dst"`
	Date   int64  `json:"date"`
	Trains string `json:"trains"`
}
