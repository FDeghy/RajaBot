package database

import (
	"sync"

	"github.com/ostafen/clover/v2"
	"gorm.io/gorm"
)

var (
	SESSION *gorm.DB
	mutex   = &sync.RWMutex{}

	cdb     *clover.DB
	rtMutex = &sync.RWMutex{}
)
