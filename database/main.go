package database

import (
	"RajaBot/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartDatabase() error {
	db, err := gorm.Open(sqlite.Open(config.Cfg.Database.Name + ".db"))
	if err != nil {
		return err
	}
	SESSION = db
	err = SESSION.AutoMigrate(
		&TgUser{},
		&TrainWR{},
	)
	if err != nil {
		return err
	}
	return nil
}

func GetTgUser(id int64) *TgUser {
	mutex.RLock()
	user := &TgUser{}
	SESSION.Where("user_id = ?", id).Take(user)
	mutex.RUnlock()
	if user.UserID != id {
		return nil
	}
	return user
}

func SaveTgUser(user *TgUser) error {
	mutex.Lock()
	res := SESSION.Create(user)
	mutex.Unlock()
	return res.Error
}

func UpdateTgUser(user *TgUser) error {
	mutex.Lock()
	res := SESSION.Save(user)
	mutex.Unlock()
	return res.Error
}

func GetTrainWR(userId int64, day int64, trainId int, src int, dst int) *TrainWR {
	mutex.RLock()
	tr := &TrainWR{}
	SESSION.Where("user_id = ? AND train_id = ? AND day = ? AND src = ? AND dst = ?", userId, trainId, day, src, dst).Take(tr)
	mutex.RUnlock()
	if tr.UserID != userId {
		return nil
	}
	return tr
}

func GetTrainWRByTid(tid uint64) *TrainWR {
	mutex.RLock()
	tr := &TrainWR{}
	SESSION.Where("id = ?", tid).Take(tr)
	mutex.RUnlock()
	if tr.Id != tid {
		return nil
	}
	return tr
}

func GetActiveTrainWRs(id int64) *[]TrainWR {
	mutex.RLock()
	tr := &[]TrainWR{}
	SESSION.Where("user_id = ? AND is_done = ?", id, false).Find(tr)
	mutex.RUnlock()
	return tr
}

func GetActiveTrainWRsByTrainId(train_id int) *[]TrainWR {
	mutex.RLock()
	tr := &[]TrainWR{}
	SESSION.Where("train_id = ? AND is_done = ?", train_id, false).Find(tr)
	mutex.RUnlock()
	return tr
}

func GetActiveTrainWRsByInfo(day int64, src int, dst int) *[]TrainWR {
	mutex.RLock()
	tr := &[]TrainWR{}
	SESSION.Where("day = ? AND src = ? AND dst = ? AND is_done = ?", day, src, dst, false).Find(tr)
	mutex.RUnlock()
	return tr
}

func SaveTrainWR(tr *TrainWR) (uint64, error) {
	mutex.Lock()
	res := SESSION.Create(tr)
	mutex.Unlock()
	return tr.Id, res.Error
}

func UpdateTrainWR(tr *TrainWR) error {
	mutex.Lock()
	res := SESSION.Save(tr)
	mutex.Unlock()
	return res.Error
}

func DeleteTrainWR(tr *TrainWR) {
	mutex.Lock()
	SESSION.Delete(tr)
	mutex.Unlock()
}
