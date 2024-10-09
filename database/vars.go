package database

import (
	"sync"

	"gorm.io/gorm"
)

var (
	SESSION *gorm.DB
	mutex   = &sync.RWMutex{}

	rtMutex = &sync.RWMutex{}
)
