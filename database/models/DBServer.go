package databases

import "gorm.io/gorm"

type DBServer struct {
	DB     *gorm.DB
	Config *DBConfig
}
