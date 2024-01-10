package database

import (
	"RajaBot/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartDatabase() error {
	conf := &gorm.Config{
		SkipDefaultTransaction: true,
	}
	db, err := gorm.Open(sqlite.Open(config.Cfg.Database.Name+".db"), conf)
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

func SaveTgUser(user *TgUser) {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Save(user)
	tx.Commit()
	mutex.Unlock()
}

func UpdateTgUser(user *TgUser) {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Save(user)
	tx.Commit()
	mutex.Unlock()
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

// by TWRID
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

func GetAllActiveTrainWRs() *[]TrainWR {
	mutex.RLock()
	tr := &[]TrainWR{}
	SESSION.Where("is_done = ? AND train_id != ?", false, 0).Find(tr)
	mutex.RUnlock()
	return tr
}

func SaveTrainWR(tr *TrainWR) uint64 {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Save(tr)
	tx.Commit()
	mutex.Unlock()
	return tr.Id
}

func UpdateTrainWR(tr *TrainWR) {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Save(tr)
	tx.Commit()
	mutex.Unlock()
}

func DeleteTrainWR(tr *TrainWR) {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Delete(tr)
	tx.Commit()
	mutex.Unlock()
}
