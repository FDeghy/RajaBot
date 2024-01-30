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
		&Subscription{},
		&Payment{},
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

func GetActiveTrainWRs(userId int64) []*TrainWR {
	mutex.RLock()
	tr := []*TrainWR{}
	SESSION.Where("user_id = ? AND is_done = ?", userId, false).Find(&tr)
	mutex.RUnlock()
	return tr
}

func FilterTrainWRsByTrainId(trainId int, trWRs []*TrainWR) []*TrainWR {
	tempTrWRs := []*TrainWR{}
	for _, tr := range trWRs {
		if tr.TrainId == trainId {
			tempTrWRs = append(tempTrWRs, tr)
		}
	}
	return tempTrWRs
}

func GetActiveTrainWRsByInfo(day int64, src int, dst int) []*TrainWR {
	mutex.RLock()
	tr := []*TrainWR{}
	SESSION.Where("day = ? AND src = ? AND dst = ? AND is_done = ?", day, src, dst, false).Find(&tr)
	mutex.RUnlock()
	return tr
}

func GetAllActiveTrainWRs() []*TrainWR {
	mutex.RLock()
	tr := []*TrainWR{}
	SESSION.Where("train_id != ? AND is_done = ?", 0, false).Find(&tr)
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

func GetSubscription(userId int64) *Subscription {
	mutex.RLock()
	sub := &Subscription{}
	SESSION.Where("user_id = ?", userId).Take(sub)
	mutex.RUnlock()
	if sub.UserID != userId {
		return nil
	}
	return sub
}

func SaveSubscription(sub *Subscription) {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Save(sub)
	tx.Commit()
	mutex.Unlock()
}

func UpdateSubscription(sub *Subscription) {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Save(sub)
	tx.Commit()
	mutex.Unlock()
}

func DeleteSubscription(sub *Subscription) {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Delete(sub)
	tx.Commit()
	mutex.Unlock()
}

func GetPayment(orderId string) *Payment {
	mutex.RLock()
	paym := &Payment{}
	SESSION.Where("order_id = ?", orderId).Take(paym)
	mutex.RUnlock()
	if paym.OrderID != orderId {
		return nil
	}
	return paym
}

func GetPaymentByTransId(transId string) *Payment {
	mutex.RLock()
	paym := &Payment{}
	SESSION.Where("trans_id = ?", transId).Take(paym)
	mutex.RUnlock()
	if paym.TransID != transId {
		return nil
	}
	return paym
}

func GetUncompletedPayment(userId int64) []*Payment {
	mutex.RLock()
	payms := []*Payment{}
	SESSION.Where("user_id = ? AND is_done = ?", userId, false).Find(&payms)
	mutex.RUnlock()
	return payms
}

func SavePayment(paym *Payment) {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Save(paym)
	tx.Commit()
	mutex.Unlock()
}

func UpdatePayment(paym *Payment) {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Save(paym)
	tx.Commit()
	mutex.Unlock()
}

func DeletePayment(paym *Payment) {
	mutex.Lock()
	tx := SESSION.Begin()
	tx.Delete(paym)
	tx.Commit()
	mutex.Unlock()
}
