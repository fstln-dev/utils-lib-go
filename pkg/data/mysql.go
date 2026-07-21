package data

import (
	"database/sql"
	"gorm.io/driver/bigquery"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGORMClient(s string, tables ...interface{}) (*gorm.DB, error) {

	sqlDB, err := sql.Open("mysql", s)
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	gorm_client, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true, // 跳过默认事务
		IgnoreRelationshipsWhenMigrating:         true,
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
		Logger:                                   logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	if len(tables) > 0 {
		if err := gorm_client.AutoMigrate(tables...); err != nil {
			return nil, err
		}
	}

	return gorm_client, nil
}

func NewBigQueryClient(s string) (*gorm.DB, error) {
	return gorm.Open(bigquery.Open(s), &gorm.Config{})
}
