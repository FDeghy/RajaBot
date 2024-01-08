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
	Src     int    `json:"src"`
	Dst     int    `json:"dst"`
	IsDone  bool   `json:"is_done"`
}
